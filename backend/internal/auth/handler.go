package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"goadmin-backend/internal/audit"
	"goadmin-backend/internal/domain"
	"goadmin-backend/internal/platform/httperr"
	"goadmin-backend/internal/platform/httpjson"
	"goadmin-backend/internal/platform/recaptcha"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type Handler struct {
	httpjson.Handler
	authService     Service
	auditService    audit.Service
	recaptchaVerifier *recaptcha.Verifier
	recaptchaMinScore float64
}

func NewHandler(authService Service, logger *slog.Logger) *Handler {
	return &Handler{
		Handler: httpjson.Handler{
			Logger: logger,
		},
		authService:      authService,
		recaptchaMinScore: 0.5,
	}
}

func (h *Handler) SetAuditService(s audit.Service) {
	h.auditService = s
}

func (h *Handler) SetRecaptcha(v *recaptcha.Verifier, minScore float64) {
	h.recaptchaVerifier = v
	h.recaptchaMinScore = minScore
}

// Login handler signs in a user.
func (h *Handler) Login(res http.ResponseWriter, req *http.Request) {
	var credentials domain.Credentials

	if err := h.ParseJSON(res, req, &credentials); err != nil {
		h.Logger.Error("error decoding credentials", slog.Any("err", err))

		return
	}

	token, err := h.authService.Login(req.Context(), credentials)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			h.Logger.Error("error invalid credentials", slog.Any("err", err))
			httperr.JSONError(res, err, http.StatusUnauthorized, req.URL.Path)
			return
		}

		if errors.Is(err, ErrAccountLocked) {
			h.Logger.Error("error account locked", slog.Any("err", err))
			httperr.JSONError(res, err, http.StatusLocked, req.URL.Path)
			return
		}

		h.Logger.Error("error logging in", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)
		return
	}

	if h.auditService != nil {
		audit.LogFromRequest(req.Context(), req, h.auditService,
			"", credentials.Username, "login", "user", "", credentials.Username)
	}

	h.RespondJSON(res, token, http.StatusOK)
}

// register handler signs up a user
func (h *Handler) Register(res http.ResponseWriter, req *http.Request) {
	var regReq RegisterRequest

	if err := h.ParseJSON(res, req, &regReq); err != nil {
		h.Logger.Error("error decoding register request", slog.Any("err", err))

		return
	}

	if h.recaptchaVerifier != nil && regReq.RecaptchaToken != "" {
		if _, err := h.recaptchaVerifier.Verify(
			regReq.RecaptchaToken, "signup", h.recaptchaMinScore,
		); err != nil {
			h.Logger.Error("recaptcha verification failed", slog.Any("err", err))
			httperr.JSONError(res, err, http.StatusBadRequest, req.URL.Path)
			return
		}
	}

	user := &domain.User{
		Username:  regReq.Username,
		Password:  regReq.Password,
		Email:     regReq.Email,
		FirstName: regReq.FirstName,
		LastName:  regReq.LastName,
	}

	newUser, err := h.authService.Register(req.Context(), user)
	if err != nil {
		h.Logger.Error("error registering user", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	res.WriteHeader(http.StatusCreated)

	if h.auditService != nil {
		audit.LogFromRequest(req.Context(), req, h.auditService,
			newUser.ID, newUser.Username, "signup", "user", newUser.ID, newUser.Username)
	}

	h.RespondJSON(res, ToUserResponse(newUser), http.StatusCreated)
}

// SignInWithGoogle handler signs in a user using Google.
func (h *Handler) SignInWithGoogle(res http.ResponseWriter, req *http.Request) {
	var idTokenReq GoogleIDTokenVerifyRequest

	if err := h.ParseJSON(res, req, &idTokenReq); err != nil {
		h.Logger.Error("error decoding id token request", slog.Any("err", err))

		return
	}

	token, err := h.authService.ValidateGoogleIDToken(
		req.Context(),
		idTokenReq.IDToken,
		"", // Use the default audience
	)
	if err != nil {
		// invalid id_token
		if errors.Is(err, ErrInvalidIDToken) {
			h.Logger.Error("error validating google id token", slog.Any("err", err))

			httperr.JSONError(res, err, http.StatusUnauthorized)

			return
		}

		// user not found
		var userNotFoundErr *domain.ResourceNotFoundError
		if errors.As(err, &userNotFoundErr) {
			h.Logger.Error("error finding user", slog.Any("err", userNotFoundErr))

			httperr.JSONError(res, userNotFoundErr, http.StatusNotFound)

			return
		}

		// other errors
		h.Logger.Error("error validating google id token", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	h.RespondJSON(res, token, http.StatusOK)
}

// Logout handler logs out the current user.
func (h *Handler) Logout(res http.ResponseWriter, req *http.Request) {
	tokenString := FindToken(req)

	if err := h.authService.Logout(req.Context(), tokenString); err != nil {
		h.Logger.Error("error logging out", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	res.WriteHeader(http.StatusOK)
}

// ForgotPassword handles password reset requests.
func (h *Handler) ForgotPassword(res http.ResponseWriter, req *http.Request) {
	var fpr ForgotPasswordRequest

	if err := h.ParseJSON(res, req, &fpr); err != nil {
		h.Logger.Error("error decoding forgot password request", slog.Any("err", err))
		return
	}

	_, err := h.authService.RequestPasswordReset(req.Context(), fpr.Email)
	if err != nil {
		h.Logger.Error("error requesting password reset", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)
		return
	}

	h.RespondJSON(res, map[string]string{
		"message": "If that email is registered, a password reset link has been sent.",
	}, http.StatusOK)
}

// ResetPassword handles password reset with a valid token.
func (h *Handler) ResetPassword(res http.ResponseWriter, req *http.Request) {
	var rpr ResetPasswordRequest

	if err := h.ParseJSON(res, req, &rpr); err != nil {
		h.Logger.Error("error decoding reset password request", slog.Any("err", err))
		return
	}

	err := h.authService.ResetPassword(req.Context(), rpr.Token, rpr.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidResetToken) || errors.Is(err, ErrResetTokenUsed) {
			h.Logger.Error("error resetting password", slog.Any("err", err))
			httperr.JSONError(res, err, http.StatusBadRequest, req.URL.Path)
			return
		}

		h.Logger.Error("error resetting password", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)
		return
	}

	h.RespondJSON(res, map[string]string{
		"message": "Password has been reset successfully.",
	}, http.StatusOK)
}

// RefreshToken handles token refresh requests.
func (h *Handler) RefreshToken(res http.ResponseWriter, req *http.Request) {
	var rtr RefreshTokenRequest

	if err := h.ParseJSON(res, req, &rtr); err != nil {
		h.Logger.Error("error decoding refresh token request", slog.Any("err", err))
		return
	}

	token, err := h.authService.RefreshToken(req.Context(), rtr.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			h.Logger.Error("error refreshing token", slog.Any("err", err))
			httperr.JSONError(res, err, http.StatusUnauthorized, req.URL.Path)
			return
		}

		h.Logger.Error("error refreshing token", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)
		return
	}

	h.RespondJSON(res, token, http.StatusOK)
}

// SetupTOTP handler generates a TOTP secret and returns a QR code URL.
func (h *Handler) SetupTOTP(res http.ResponseWriter, req *http.Request) {
	tokenString := FindToken(req)

	secret, qrURL, err := h.authService.SetupTOTP(req.Context(), tokenString)
	if err != nil {
		h.Logger.Error("error setting up totp", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)
		return
	}

	h.RespondJSON(res, map[string]string{
		"secret":  secret,
		"qr_url":  qrURL,
		"message": "Scan the QR code with your authenticator app, then verify with a code.",
	}, http.StatusOK)
}

// VerifyTOTP handler verifies a TOTP code and enables 2FA.
func (h *Handler) VerifyTOTP(res http.ResponseWriter, req *http.Request) {
	tokenString := FindToken(req)

	var body struct {
		Code string `json:"code"`
	}
	if err := h.ParseJSON(res, req, &body); err != nil {
		return
	}

	if err := h.authService.VerifyTOTP(req.Context(), tokenString, body.Code); err != nil {
		h.Logger.Error("error verifying totp", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusBadRequest, req.URL.Path)
		return
	}

	h.RespondJSON(res, map[string]string{
		"message": "2FA has been enabled successfully.",
	}, http.StatusOK)
}

// DisableTOTP handler disables 2FA for the current user.
func (h *Handler) DisableTOTP(res http.ResponseWriter, req *http.Request) {
	tokenString := FindToken(req)

	var body struct {
		Password string `json:"password"`
	}
	if err := h.ParseJSON(res, req, &body); err != nil {
		return
	}

	if err := h.authService.DisableTOTP(req.Context(), tokenString, body.Password); err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httperr.JSONError(res, err, http.StatusUnauthorized, req.URL.Path)
			return
		}
		h.Logger.Error("error disabling totp", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)
		return
	}

	h.RespondJSON(res, map[string]string{
		"message": "2FA has been disabled.",
	}, http.StatusOK)
}

// Login2FA handler completes login with a TOTP code.
func (h *Handler) Login2FA(res http.ResponseWriter, req *http.Request) {
	var body struct {
		TempToken string `json:"temp_token"`
		Code      string `json:"code"`
	}
	if err := h.ParseJSON(res, req, &body); err != nil {
		return
	}

	token, err := h.authService.Login2FA(req.Context(), body.TempToken, body.Code)
	if err != nil {
		h.Logger.Error("error 2fa login", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusUnauthorized, req.URL.Path)
		return
	}

	h.RespondJSON(res, token, http.StatusOK)
}

// Profile handler gets the profile of the current user.
func (h *Handler) Profile(res http.ResponseWriter, req *http.Request) {
	tokenString := FindToken(req)

	user, err := h.authService.Profile(req.Context(), tokenString)
	if err != nil {
		h.Logger.Error("error getting profile", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	h.RespondJSON(res, ToUserResponse(user), http.StatusOK)
}

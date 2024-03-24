package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"goadmin-backend/internal/domain"
	"goadmin-backend/internal/platform/httperr"
	"goadmin-backend/internal/platform/httpjson"
)

type Handler struct {
	httpjson.Handler
	authService Service
}

func NewHandler(authService Service, logger *slog.Logger) *Handler {
	return &Handler{
		Handler: httpjson.Handler{
			Logger: logger,
		},
		authService: authService,
	}
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

		h.Logger.Error("error logging in", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
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

// Logout handler logs out a user.
func (h *Handler) Logout(res http.ResponseWriter, req *http.Request) {
	tokenString := FindToken(req)

	if err := h.authService.Logout(req.Context(), tokenString); err != nil {
		h.Logger.Error("error logging out", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	res.WriteHeader(http.StatusOK)
}

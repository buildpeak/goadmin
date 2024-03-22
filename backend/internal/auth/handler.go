package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"goadmin-backend/internal/domain"
	"goadmin-backend/internal/platform/httperr"
)

type Handler struct {
	authService Service
	logger      *slog.Logger
}

func NewHandler(authService Service, logger *slog.Logger) *Handler {
	return &Handler{authService: authService, logger: logger}
}

// Login handler signs in a user.
func (h *Handler) Login(res http.ResponseWriter, req *http.Request) {
	var credentials domain.Credentials

	if err := json.NewDecoder(req.Body).Decode(&credentials); err != nil {
		h.logger.Error("error decoding credentials", slog.Any("err", err))

		http.Error(res, "invalid request", http.StatusBadRequest)

		return
	}

	token, err := h.authService.Login(req.Context(), credentials)
	if err != nil {
		h.logger.Error("error logging in", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	if err := json.NewEncoder(res).Encode(token); err != nil {
		h.logger.Error("error writing token", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}
}

// register handler signs up a user
func (h *Handler) Register(res http.ResponseWriter, req *http.Request) {
	var user domain.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		h.logger.Error("error decoding user", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	newUser, err := h.authService.Register(req.Context(), &user)
	if err != nil {
		h.logger.Error("error registering user", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	res.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(res).Encode(ToUserResponse(newUser)); err != nil {
		h.logger.Error("error encoding user", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}
}

// SignInWithGoogle handler signs in a user using Google.
func (h *Handler) SignInWithGoogle(res http.ResponseWriter, req *http.Request) {
	var idTokenReq GoogleIDTokenVerifyRequest

	if err := json.NewDecoder(req.Body).Decode(&idTokenReq); err != nil {
		h.logger.Error("error decoding id token request", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

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
			h.logger.Error("error validating google id token", slog.Any("err", err))

			httperr.JSONError(res, err, http.StatusUnauthorized)

			return
		}

		// user not found
		var userNotFoundErr *domain.ResourceNotFoundError
		if errors.As(err, &userNotFoundErr) {
			h.logger.Error("error finding user", slog.Any("err", userNotFoundErr))

			httperr.JSONError(res, userNotFoundErr, http.StatusNotFound)

			return
		}

		// other errors
		h.logger.Error("error validating google id token", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	if err := json.NewEncoder(res).Encode(token); err != nil {
		h.logger.Error("error encoding user", slog.Any("err", err))

		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}
}

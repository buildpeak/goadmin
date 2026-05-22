package user

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"goadmin-backend/internal/auth"
	"goadmin-backend/internal/domain"
	"goadmin-backend/internal/platform/httperr"
	"goadmin-backend/internal/platform/httpjson"
)

type Handler struct {
	httpjson.Handler
	userService Service
}

func NewHandler(userService Service, logger *slog.Logger) *Handler {
	return &Handler{
		Handler: httpjson.Handler{
			Logger: logger,
		},
		userService: userService,
	}
}

func (h *Handler) List(res http.ResponseWriter, req *http.Request) {
	filter := &domain.UserFilter{}

	users, err := h.userService.List(req.Context(), filter)
	if err != nil {
		h.Logger.Error("error listing users", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, users, http.StatusOK)
}

func (h *Handler) GetByID(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	user, err := h.userService.GetByID(req.Context(), id)
	if err != nil {
		h.Logger.Error("error getting user", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusNotFound, req.URL.Path)

		return
	}

	h.RespondJSON(res, auth.ToUserResponse(user), http.StatusOK)
}

func (h *Handler) Update(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	var updates domain.User
	if err := h.ParseJSON(res, req, &updates); err != nil {
		h.Logger.Error("error decoding update request", slog.Any("err", err))

		return
	}

	updates.ID = id

	updated, err := h.userService.Update(req.Context(), &updates)
	if err != nil {
		h.Logger.Error("error updating user", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, auth.ToUserResponse(updated), http.StatusOK)
}

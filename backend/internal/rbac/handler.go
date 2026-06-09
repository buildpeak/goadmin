package rbac

import (
	"errors"
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
	rbacService Service
}

func NewHandler(rbacService Service, logger *slog.Logger) *Handler {
	return &Handler{
		Handler:     httpjson.Handler{Logger: logger},
		rbacService: rbacService,
	}
}

// Role handlers

func (h *Handler) ListRoles(res http.ResponseWriter, req *http.Request) {
	roles, err := h.rbacService.ListRoles(req.Context())
	if err != nil {
		h.Logger.Error("error listing roles", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, roles, http.StatusOK)
}

func (h *Handler) CreateRole(res http.ResponseWriter, req *http.Request) {
	var body struct {
		Name string `json:"name"`
	}

	if err := h.ParseJSON(res, req, &body); err != nil {
		return
	}

	role, err := h.rbacService.CreateRole(req.Context(), body.Name)
	if err != nil {
		h.Logger.Error("error creating role", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, role, http.StatusCreated)
}

func (h *Handler) GetRole(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	role, err := h.rbacService.GetRole(req.Context(), id)
	if err != nil {
		h.Logger.Error("error getting role", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, role, http.StatusOK)
}

func (h *Handler) UpdateRole(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	var body domain.Role
	if err := h.ParseJSON(res, req, &body); err != nil {
		return
	}

	body.ID = id

	role, err := h.rbacService.UpdateRole(req.Context(), &body)
	if err != nil {
		h.Logger.Error("error updating role", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, role, http.StatusOK)
}

func (h *Handler) DeleteRole(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	if err := h.rbacService.DeleteRole(req.Context(), id); err != nil {
		h.Logger.Error("error deleting role", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	res.WriteHeader(http.StatusNoContent)
}

// Permission handlers

func (h *Handler) ListPermissions(res http.ResponseWriter, req *http.Request) {
	perms, err := h.rbacService.ListPermissions(req.Context())
	if err != nil {
		h.Logger.Error("error listing permissions", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, perms, http.StatusOK)
}

func (h *Handler) CreatePermission(res http.ResponseWriter, req *http.Request) {
	var body struct {
		Name     string `json:"name"`
		RuleType string `json:"rule_type"`
		Rule     string `json:"rule"`
	}

	if err := h.ParseJSON(res, req, &body); err != nil {
		return
	}

	perm, err := h.rbacService.CreatePermission(req.Context(), body.Name, body.RuleType, body.Rule)
	if err != nil {
		h.Logger.Error("error creating permission", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, perm, http.StatusCreated)
}

func (h *Handler) GetPermission(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	perm, err := h.rbacService.GetPermission(req.Context(), id)
	if err != nil {
		h.Logger.Error("error getting permission", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, perm, http.StatusOK)
}

func (h *Handler) UpdatePermission(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	var body domain.AuthPermission
	if err := h.ParseJSON(res, req, &body); err != nil {
		return
	}

	body.ID = id

	perm, err := h.rbacService.UpdatePermission(req.Context(), &body)
	if err != nil {
		h.Logger.Error("error updating permission", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, perm, http.StatusOK)
}

func (h *Handler) DeletePermission(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	if err := h.rbacService.DeletePermission(req.Context(), id); err != nil {
		h.Logger.Error("error deleting permission", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	res.WriteHeader(http.StatusNoContent)
}

// Role-Permission assignment handlers

func (h *Handler) ListPermissionsByRole(res http.ResponseWriter, req *http.Request) {
	roleID := chi.URLParam(req, "id")

	perms, err := h.rbacService.ListPermissionsByRole(req.Context(), roleID)
	if err != nil {
		h.Logger.Error("error listing permissions by role", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, perms, http.StatusOK)
}

func (h *Handler) AddPermissionToRole(res http.ResponseWriter, req *http.Request) {
	roleID := chi.URLParam(req, "id")
	permissionID := chi.URLParam(req, "permissionId")

	if err := h.rbacService.AssignPermissionToRole(req.Context(), roleID, permissionID); err != nil {
		h.Logger.Error("error assigning permission to role", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemovePermissionFromRole(res http.ResponseWriter, req *http.Request) {
	roleID := chi.URLParam(req, "id")
	permissionID := chi.URLParam(req, "permissionId")

	if err := h.rbacService.UnassignPermissionFromRole(req.Context(), roleID, permissionID); err != nil {
		h.Logger.Error("error removing permission from role", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	res.WriteHeader(http.StatusNoContent)
}

// User-Role assignment handlers

func (h *Handler) ListRolesByUser(res http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "userId")

	roles, err := h.rbacService.ListRolesByUser(req.Context(), userID)
	if err != nil {
		h.Logger.Error("error listing roles by user", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, roles, http.StatusOK)
}

func (h *Handler) AssignRoleToUser(res http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "userId")
	roleID := chi.URLParam(req, "roleId")

	if err := h.rbacService.AssignRoleToUser(req.Context(), userID, roleID); err != nil {
		h.Logger.Error("error assigning role to user", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UnassignRoleFromUser(res http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "userId")
	roleID := chi.URLParam(req, "roleId")

	if err := h.rbacService.UnassignRoleFromUser(req.Context(), userID, roleID); err != nil {
		h.Logger.Error("error unassigning role from user", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	res.WriteHeader(http.StatusNoContent)
}

// User-Permission assignment handlers

func (h *Handler) ListPermissionsByUser(res http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "userId")

	perms, err := h.rbacService.ListPermissionsByUser(req.Context(), userID)
	if err != nil {
		h.Logger.Error("error listing permissions by user", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	h.RespondJSON(res, perms, http.StatusOK)
}

func (h *Handler) AssignPermissionToUser(res http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "userId")
	permissionID := chi.URLParam(req, "permissionId")

	if err := h.rbacService.AssignPermissionToUser(req.Context(), userID, permissionID); err != nil {
		h.Logger.Error("error assigning permission to user", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UnassignPermissionFromUser(res http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "userId")
	permissionID := chi.URLParam(req, "permissionId")

	if err := h.rbacService.UnassignPermissionFromUser(req.Context(), userID, permissionID); err != nil {
		h.Logger.Error("error unassigning permission from user", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

		return
	}

	res.WriteHeader(http.StatusNoContent)
}

// RequirePermission returns middleware that checks the authenticated user
// has the named permission. Must be used after auth.Authenticator().
func (h *Handler) RequirePermission(permissionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			user, ok := auth.UserFromContext(req.Context())
			if !ok {
				httperr.JSONError(res, errors.New("unauthorized"), http.StatusUnauthorized, req.URL.Path)

				return
			}

			has, err := h.rbacService.HasPermission(req.Context(), user.ID, permissionName)
			if err != nil {
				h.Logger.Error("error checking permission",
					slog.String("permission", permissionName),
					slog.Any("err", err),
				)
				httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)

				return
			}

			if !has {
				httperr.JSONError(res, errors.New("forbidden"), http.StatusForbidden, req.URL.Path)

				return
			}

			next.ServeHTTP(res, req)
		})
	}
}

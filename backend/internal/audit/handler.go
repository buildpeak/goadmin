package audit

import (
	"log/slog"
	"net/http"
	"strconv"

	"goadmin-backend/internal/domain"
	"goadmin-backend/internal/platform/httperr"
	"goadmin-backend/internal/platform/httpjson"
)

type Handler struct {
	httpjson.Handler
	auditService Service
}

func NewHandler(auditService Service, logger *slog.Logger) *Handler {
	return &Handler{
		Handler: httpjson.Handler{Logger: logger},
		auditService: auditService,
	}
}

func (h *Handler) List(res http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	filter := &domain.AuditLogFilter{
		ActorUsername: q.Get("actor"),
		Action:        q.Get("action"),
		ResourceType:  q.Get("resource"),
		Limit:         limit,
		Offset:        offset,
	}

	logs, err := h.auditService.List(req.Context(), filter)
	if err != nil {
		h.Logger.Error("error listing audit logs", slog.Any("err", err))
		httperr.JSONError(res, err, http.StatusInternalServerError, req.URL.Path)
		return
	}

	if logs == nil {
		logs = []domain.AuditLog{}
	}

	h.RespondJSON(res, logs, http.StatusOK)
}

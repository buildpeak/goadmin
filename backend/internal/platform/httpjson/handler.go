package httpjson

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"goadmin-backend/internal/platform/httperr"
)

type Handler struct {
	logger *slog.Logger
}

func (h *Handler) ParseJSON(
	res http.ResponseWriter,
	req *http.Request,
	dataPtr any,
) error {
	// Decode the request body into the dataPtr.
	if err := json.NewDecoder(req.Body).Decode(dataPtr); err != nil {
		httperr.JSONError(res, err, http.StatusBadRequest)

		return err //nolint:wrapcheck // no need to wrap
	}

	return nil
}

func (h *Handler) RespondJSON(res http.ResponseWriter, data any, status int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)

	if data == nil {
		return
	}

	if err := json.NewEncoder(res).Encode(data); err != nil {
		h.logger.Error("error encoding response", slog.Any("err", err))
	}
}

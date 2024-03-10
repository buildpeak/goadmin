package user

import "net/http"

type Handler struct {
	userService Service
}

func NewHandler(userService Service) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) List(_ http.ResponseWriter, _ *http.Request) {
	// ...
}

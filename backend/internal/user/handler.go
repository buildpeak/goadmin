package user

import "net/http"

type Handler struct {
	userService UserService
}

func NewHandler(userService UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// ...
}

package auth

type Handler struct {
	authService Service
}

func NewHandler(authService Service) *Handler {
	return &Handler{authService: authService}
}

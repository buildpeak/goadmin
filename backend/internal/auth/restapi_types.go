package auth

// LoginRequest represents a request to sign in a user.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a response to sign in a user.
type LoginResponse struct {
	Token string `json:"token"`
}

// LogoutRequest represents a request to invalidate a token.
type LogoutRequest struct {
	Token string `json:"token" validate:"required"`
}

// RegisterRequest represents a request to register a user.
type RegisterRequest struct {
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

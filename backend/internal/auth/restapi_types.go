package auth

import (
	"time"

	"goadmin-backend/internal/domain"
)

// LoginRequest represents a request to sign in a user.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
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

type GoogleIDTokenVerifyRequest struct {
	IDToken    string `json:"id_token" validate:"required"`
	GCSRFToken string `json:"g_csrf_token"`
}

type UserResponse struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Active    bool       `json:"active"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func ToUserResponse(usr *domain.User) UserResponse {
	return UserResponse{
		ID:        usr.ID,
		Username:  usr.Username,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Email:     usr.Email,
		Active:    usr.Active,
		DeletedAt: usr.DeletedAt,
	}
}

package domain

import (
	"context"
	"time"
)

// User model
// User model defines the structure of a user
type User struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Active      bool       `json:"active"`
	Picture     string     `json:"picture"`
	TotpSecret          *string    `json:"totp_secret"`
	TotpEnabled         bool       `json:"totp_enabled"`
	FailedLoginAttempts int        `json:"failed_login_attempts"`
	LockedUntil         *time.Time `json:"locked_until"`
	Deleted             bool       `json:"deleted"`
	DeletedAt   *time.Time `json:"deleted_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type UserFilter struct {
	Email          string       `json:"email"`
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	Active         *bool        `json:"active"`
	Deleted        *bool        `json:"deleted"`
	CreatedBetween [2]time.Time `json:"created_between"`
}

type UserRole struct {
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository defines the methods that a user repository should implement
type UserRepository interface {
	FindAll(ctx context.Context, filter *UserFilter) ([]User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	UpdatePassword(ctx context.Context, id string, hashedPassword string) error
	UpdateFailedAttempts(ctx context.Context, id string, attempts int, lockedUntil *time.Time) error
	SoftDelete(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

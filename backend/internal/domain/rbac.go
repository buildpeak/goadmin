package domain

import (
	"context"
	"time"
)

type Role struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthPermission struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	RuleType  string    `json:"rule_type,omitempty"`
	Rule      string    `json:"rule,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleRepository interface {
	FindAll(ctx context.Context) ([]Role, error)
	FindByID(ctx context.Context, id string) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	Create(ctx context.Context, role *Role) (*Role, error)
	Update(ctx context.Context, role *Role) (*Role, error)
	Delete(ctx context.Context, id string) error

	FindRolesByUser(ctx context.Context, userID string) ([]Role, error)
	AssignToUser(ctx context.Context, userID, roleID string) error
	UnassignFromUser(ctx context.Context, userID, roleID string) error

	FindPermissionsByRole(ctx context.Context, roleID string) ([]AuthPermission, error)
	AddPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error
}

type PermissionRepository interface {
	FindAll(ctx context.Context) ([]AuthPermission, error)
	FindByID(ctx context.Context, id string) (*AuthPermission, error)
	FindByName(ctx context.Context, name string) (*AuthPermission, error)
	Create(ctx context.Context, permission *AuthPermission) (*AuthPermission, error)
	Update(ctx context.Context, permission *AuthPermission) (*AuthPermission, error)
	Delete(ctx context.Context, id string) error

	FindPermissionsByUser(ctx context.Context, userID string) ([]AuthPermission, error)
	AssignToUser(ctx context.Context, userID, permissionID string) error
	UnassignFromUser(ctx context.Context, userID, permissionID string) error

	FindEffectivePermissions(ctx context.Context, userID string) ([]AuthPermission, error)
}

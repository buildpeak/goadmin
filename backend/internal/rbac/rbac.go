package rbac

import (
	"context"
	"fmt"

	"goadmin-backend/internal/domain"
)

type Service interface {
	// Role management
	CreateRole(ctx context.Context, name string) (*domain.Role, error)
	GetRole(ctx context.Context, id string) (*domain.Role, error)
	ListRoles(ctx context.Context) ([]domain.Role, error)
	UpdateRole(ctx context.Context, role *domain.Role) (*domain.Role, error)
	DeleteRole(ctx context.Context, id string) error

	// Permission management
	CreatePermission(ctx context.Context, name, ruleType, rule string) (*domain.AuthPermission, error)
	GetPermission(ctx context.Context, id string) (*domain.AuthPermission, error)
	ListPermissions(ctx context.Context) ([]domain.AuthPermission, error)
	UpdatePermission(ctx context.Context, perm *domain.AuthPermission) (*domain.AuthPermission, error)
	DeletePermission(ctx context.Context, id string) error

	// Role-Permission assignments
	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error
	UnassignPermissionFromRole(ctx context.Context, roleID, permissionID string) error
	ListPermissionsByRole(ctx context.Context, roleID string) ([]domain.AuthPermission, error)

	// User-Role assignments
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	UnassignRoleFromUser(ctx context.Context, userID, roleID string) error
	ListRolesByUser(ctx context.Context, userID string) ([]domain.Role, error)

	// User-Permission assignments (direct)
	AssignPermissionToUser(ctx context.Context, userID, permissionID string) error
	UnassignPermissionFromUser(ctx context.Context, userID, permissionID string) error
	ListPermissionsByUser(ctx context.Context, userID string) ([]domain.AuthPermission, error)

	// Permission checks
	HasPermission(ctx context.Context, userID, permissionName string) (bool, error)
	GetUserPermissions(ctx context.Context, userID string) ([]domain.AuthPermission, error)
}

type rbacService struct {
	roleRepo       domain.RoleRepository
	permissionRepo domain.PermissionRepository
}

var _ Service = &rbacService{}

func NewService(roleRepo domain.RoleRepository, permRepo domain.PermissionRepository) Service {
	return &rbacService{
		roleRepo:       roleRepo,
		permissionRepo: permRepo,
	}
}

// Role management

func (s *rbacService) CreateRole(ctx context.Context, name string) (*domain.Role, error) {
	role := &domain.Role{Name: name}
	created, err := s.roleRepo.Create(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}

	return created, nil
}

func (s *rbacService) GetRole(ctx context.Context, id string) (*domain.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get role: %w", err)
	}

	return role, nil
}

func (s *rbacService) ListRoles(ctx context.Context) ([]domain.Role, error) {
	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}

	return roles, nil
}

func (s *rbacService) UpdateRole(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	updated, err := s.roleRepo.Update(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}

	return updated, nil
}

func (s *rbacService) DeleteRole(ctx context.Context, id string) error {
	if err := s.roleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}

	return nil
}

// Permission management

func (s *rbacService) CreatePermission(ctx context.Context, name, ruleType, rule string) (*domain.AuthPermission, error) {
	perm := &domain.AuthPermission{Name: name, RuleType: ruleType, Rule: rule}
	created, err := s.permissionRepo.Create(ctx, perm)
	if err != nil {
		return nil, fmt.Errorf("create permission: %w", err)
	}

	return created, nil
}

func (s *rbacService) GetPermission(ctx context.Context, id string) (*domain.AuthPermission, error) {
	perm, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get permission: %w", err)
	}

	return perm, nil
}

func (s *rbacService) ListPermissions(ctx context.Context) ([]domain.AuthPermission, error) {
	perms, err := s.permissionRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}

	return perms, nil
}

func (s *rbacService) UpdatePermission(ctx context.Context, perm *domain.AuthPermission) (*domain.AuthPermission, error) {
	updated, err := s.permissionRepo.Update(ctx, perm)
	if err != nil {
		return nil, fmt.Errorf("update permission: %w", err)
	}

	return updated, nil
}

func (s *rbacService) DeletePermission(ctx context.Context, id string) error {
	if err := s.permissionRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete permission: %w", err)
	}

	return nil
}

// Role-Permission assignments

func (s *rbacService) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	if err := s.roleRepo.AddPermissionToRole(ctx, roleID, permissionID); err != nil {
		return fmt.Errorf("assign permission to role: %w", err)
	}

	return nil
}

func (s *rbacService) UnassignPermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	if err := s.roleRepo.RemovePermissionFromRole(ctx, roleID, permissionID); err != nil {
		return fmt.Errorf("unassign permission from role: %w", err)
	}

	return nil
}

func (s *rbacService) ListPermissionsByRole(ctx context.Context, roleID string) ([]domain.AuthPermission, error) {
	perms, err := s.roleRepo.FindPermissionsByRole(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("list permissions by role: %w", err)
	}

	return perms, nil
}

// User-Role assignments

func (s *rbacService) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	if err := s.roleRepo.AssignToUser(ctx, userID, roleID); err != nil {
		return fmt.Errorf("assign role to user: %w", err)
	}

	return nil
}

func (s *rbacService) UnassignRoleFromUser(ctx context.Context, userID, roleID string) error {
	if err := s.roleRepo.UnassignFromUser(ctx, userID, roleID); err != nil {
		return fmt.Errorf("unassign role from user: %w", err)
	}

	return nil
}

func (s *rbacService) ListRolesByUser(ctx context.Context, userID string) ([]domain.Role, error) {
	roles, err := s.roleRepo.FindRolesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list roles by user: %w", err)
	}

	return roles, nil
}

// User-Permission assignments

func (s *rbacService) AssignPermissionToUser(ctx context.Context, userID, permissionID string) error {
	if err := s.permissionRepo.AssignToUser(ctx, userID, permissionID); err != nil {
		return fmt.Errorf("assign permission to user: %w", err)
	}

	return nil
}

func (s *rbacService) UnassignPermissionFromUser(ctx context.Context, userID, permissionID string) error {
	if err := s.permissionRepo.UnassignFromUser(ctx, userID, permissionID); err != nil {
		return fmt.Errorf("unassign permission from user: %w", err)
	}

	return nil
}

func (s *rbacService) ListPermissionsByUser(ctx context.Context, userID string) ([]domain.AuthPermission, error) {
	perms, err := s.permissionRepo.FindPermissionsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list permissions by user: %w", err)
	}

	return perms, nil
}

// Permission checks

func (s *rbacService) HasPermission(ctx context.Context, userID, permissionName string) (bool, error) {
	perms, err := s.permissionRepo.FindEffectivePermissions(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("check permission: %w", err)
	}

	for _, p := range perms {
		if p.Name == permissionName {
			return true, nil
		}
	}

	return false, nil
}

func (s *rbacService) GetUserPermissions(ctx context.Context, userID string) ([]domain.AuthPermission, error) {
	perms, err := s.permissionRepo.FindEffectivePermissions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user permissions: %w", err)
	}

	return perms, nil
}

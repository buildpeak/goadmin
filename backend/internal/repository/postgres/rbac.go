package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"goadmin-backend/internal/domain"
)

// RoleRepo

var _ domain.RoleRepository = &RoleRepo{}

type RoleRepo struct {
	db Queryer
}

func NewRoleRepo(db Queryer) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) FindAll(ctx context.Context) ([]domain.Role, error) {
	q := fmt.Sprintf(`SELECT * FROM %s ORDER BY id`, roleTable)
	results, err := query[domain.Role](ctx, r.db, q)
	if err != nil {
		return nil, fmt.Errorf("find all roles: %w", err)
	}

	roles := make([]domain.Role, len(results))
	for i, rr := range results {
		roles[i] = *rr
	}

	return roles, nil
}

func (r *RoleRepo) FindByID(ctx context.Context, id string) (*domain.Role, error) {
	q := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, roleTable)
	role, err := queryRow[domain.Role](ctx, r.db, q, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("Role", "id="+id)
		}
		return nil, fmt.Errorf("find role by ID: %w", err)
	}

	return role, nil
}

func (r *RoleRepo) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	q := fmt.Sprintf(`SELECT * FROM %s WHERE name = $1`, roleTable)
	role, err := queryRow[domain.Role](ctx, r.db, q, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("Role", "name="+name)
		}
		return nil, fmt.Errorf("find role by name: %w", err)
	}

	return role, nil
}

func (r *RoleRepo) Create(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	q := fmt.Sprintf(
		`INSERT INTO %s (name) VALUES ($1) RETURNING *`,
		roleTable,
	)
	result, err := queryRow[domain.Role](ctx, r.db, q, role.Name)
	if err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}

	return result, nil
}

func (r *RoleRepo) Update(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	q := fmt.Sprintf(
		`UPDATE %s SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING *`,
		roleTable,
	)
	result, err := queryRow[domain.Role](ctx, r.db, q, role.Name, role.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("Role", "id="+role.ID)
		}
		return nil, fmt.Errorf("update role: %w", err)
	}

	return result, nil
}

func (r *RoleRepo) Delete(ctx context.Context, id string) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, roleTable)
	_, err := exec(ctx, r.db, q, id)
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}

	return nil
}

func (r *RoleRepo) FindRolesByUser(ctx context.Context, userID string) ([]domain.Role, error) {
	q := fmt.Sprintf(
		`SELECT r.* FROM %s r JOIN %s ur ON ur.role_id = r.id WHERE ur.user_id = $1 ORDER BY r.id`,
		roleTable, userRoleTable,
	)
	results, err := query[domain.Role](ctx, r.db, q, userID)
	if err != nil {
		return nil, fmt.Errorf("find roles by user: %w", err)
	}

	roles := make([]domain.Role, len(results))
	for i, rr := range results {
		roles[i] = *rr
	}

	return roles, nil
}

func (r *RoleRepo) AssignToUser(ctx context.Context, userID, roleID string) error {
	q := fmt.Sprintf(
		`INSERT INTO %s (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userRoleTable,
	)
	_, err := exec(ctx, r.db, q, userID, roleID)
	if err != nil {
		return fmt.Errorf("assign role to user: %w", err)
	}

	return nil
}

func (r *RoleRepo) UnassignFromUser(ctx context.Context, userID, roleID string) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1 AND role_id = $2`, userRoleTable)
	_, err := exec(ctx, r.db, q, userID, roleID)
	if err != nil {
		return fmt.Errorf("unassign role from user: %w", err)
	}

	return nil
}

func (r *RoleRepo) FindPermissionsByRole(ctx context.Context, roleID string) ([]domain.AuthPermission, error) {
	q := fmt.Sprintf(
		`SELECT p.* FROM %s p JOIN %s rp ON rp.permission_id = p.id WHERE rp.role_id = $1 ORDER BY p.id`,
		permissionTable, rolePermissionTable,
	)
	results, err := query[domain.AuthPermission](ctx, r.db, q, roleID)
	if err != nil {
		return nil, fmt.Errorf("find permissions by role: %w", err)
	}

	perms := make([]domain.AuthPermission, len(results))
	for i, rp := range results {
		perms[i] = *rp
	}

	return perms, nil
}

func (r *RoleRepo) AddPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	q := fmt.Sprintf(
		`INSERT INTO %s (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		rolePermissionTable,
	)
	_, err := exec(ctx, r.db, q, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("add permission to role: %w", err)
	}

	return nil
}

func (r *RoleRepo) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	q := fmt.Sprintf(
		`DELETE FROM %s WHERE role_id = $1 AND permission_id = $2`,
		rolePermissionTable,
	)
	_, err := exec(ctx, r.db, q, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("remove permission from role: %w", err)
	}

	return nil
}

// PermissionRepo

var _ domain.PermissionRepository = &PermissionRepo{}

type PermissionRepo struct {
	db Queryer
}

func NewPermissionRepo(db Queryer) *PermissionRepo {
	return &PermissionRepo{db: db}
}

func (r *PermissionRepo) FindAll(ctx context.Context) ([]domain.AuthPermission, error) {
	q := fmt.Sprintf(`SELECT * FROM %s ORDER BY id`, permissionTable)
	results, err := query[domain.AuthPermission](ctx, r.db, q)
	if err != nil {
		return nil, fmt.Errorf("find all permissions: %w", err)
	}

	perms := make([]domain.AuthPermission, len(results))
	for i, rp := range results {
		perms[i] = *rp
	}

	return perms, nil
}

func (r *PermissionRepo) FindByID(ctx context.Context, id string) (*domain.AuthPermission, error) {
	q := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, permissionTable)
	perm, err := queryRow[domain.AuthPermission](ctx, r.db, q, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("Permission", "id="+id)
		}
		return nil, fmt.Errorf("find permission by ID: %w", err)
	}

	return perm, nil
}

func (r *PermissionRepo) FindByName(ctx context.Context, name string) (*domain.AuthPermission, error) {
	q := fmt.Sprintf(`SELECT * FROM %s WHERE name = $1`, permissionTable)
	perm, err := queryRow[domain.AuthPermission](ctx, r.db, q, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("Permission", "name="+name)
		}
		return nil, fmt.Errorf("find permission by name: %w", err)
	}

	return perm, nil
}

func (r *PermissionRepo) Create(ctx context.Context, perm *domain.AuthPermission) (*domain.AuthPermission, error) {
	q := fmt.Sprintf(
		`INSERT INTO %s (name, rule_type, rule) VALUES ($1, $2, $3) RETURNING *`,
		permissionTable,
	)
	result, err := queryRow[domain.AuthPermission](ctx, r.db, q, perm.Name, perm.RuleType, perm.Rule)
	if err != nil {
		return nil, fmt.Errorf("create permission: %w", err)
	}

	return result, nil
}

func (r *PermissionRepo) Update(ctx context.Context, perm *domain.AuthPermission) (*domain.AuthPermission, error) {
	q := fmt.Sprintf(
		`UPDATE %s SET name = $1, rule_type = $2, rule = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4 RETURNING *`,
		permissionTable,
	)
	result, err := queryRow[domain.AuthPermission](ctx, r.db, q, perm.Name, perm.RuleType, perm.Rule, perm.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("Permission", "id="+perm.ID)
		}
		return nil, fmt.Errorf("update permission: %w", err)
	}

	return result, nil
}

func (r *PermissionRepo) Delete(ctx context.Context, id string) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, permissionTable)
	_, err := exec(ctx, r.db, q, id)
	if err != nil {
		return fmt.Errorf("delete permission: %w", err)
	}

	return nil
}

func (r *PermissionRepo) FindPermissionsByUser(ctx context.Context, userID string) ([]domain.AuthPermission, error) {
	q := fmt.Sprintf(
		`SELECT p.* FROM %s p JOIN %s up ON up.permission_id = p.id WHERE up.user_id = $1 ORDER BY p.id`,
		permissionTable, userPermissionTable,
	)
	results, err := query[domain.AuthPermission](ctx, r.db, q, userID)
	if err != nil {
		return nil, fmt.Errorf("find permissions by user: %w", err)
	}

	perms := make([]domain.AuthPermission, len(results))
	for i, rp := range results {
		perms[i] = *rp
	}

	return perms, nil
}

func (r *PermissionRepo) AssignToUser(ctx context.Context, userID, permissionID string) error {
	q := fmt.Sprintf(
		`INSERT INTO %s (user_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userPermissionTable,
	)
	_, err := exec(ctx, r.db, q, userID, permissionID)
	if err != nil {
		return fmt.Errorf("assign permission to user: %w", err)
	}

	return nil
}

func (r *PermissionRepo) UnassignFromUser(ctx context.Context, userID, permissionID string) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1 AND permission_id = $2`, userPermissionTable)
	_, err := exec(ctx, r.db, q, userID, permissionID)
	if err != nil {
		return fmt.Errorf("unassign permission from user: %w", err)
	}

	return nil
}

func (r *PermissionRepo) FindEffectivePermissions(ctx context.Context, userID string) ([]domain.AuthPermission, error) {
	q := fmt.Sprintf(
		`SELECT DISTINCT p.* FROM %s p
		JOIN %s rp ON rp.permission_id = p.id
		JOIN %s ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = $1
		UNION
		SELECT DISTINCT p.* FROM %s p
		JOIN %s up ON up.permission_id = p.id
		WHERE up.user_id = $1
		ORDER BY id`,
		permissionTable, rolePermissionTable, userRoleTable,
		permissionTable, userPermissionTable,
	)
	results, err := query[domain.AuthPermission](ctx, r.db, q, userID)
	if err != nil {
		return nil, fmt.Errorf("find effective permissions: %w", err)
	}

	perms := make([]domain.AuthPermission, len(results))
	for i, rp := range results {
		perms[i] = *rp
	}

	return perms, nil
}

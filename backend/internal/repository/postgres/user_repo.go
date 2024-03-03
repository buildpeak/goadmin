package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"goadmin-backend/internal/domain"
)

var _ domain.UserRepository = &UserRepo{}

const (
	userTable = `"user"`
)

type UserRepo struct {
	db Queryer
}

func NewUserRepo(db Queryer) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// FindAll returns all users from the database
func (r *UserRepo) FindAll(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error) {
	where := "1 = 1"
	args := []interface{}{}

	switch {
	case filter.FirstName != "":
		where += " AND first_name = $1"

		args = append(args, filter.FirstName)
	case filter.LastName != "":
		where += " AND last_name = $2"

		args = append(args, filter.LastName)
	case filter.Active != nil:
		where += " AND active = $3"

		args = append(args, filter.Active)
	case filter.Deleted != nil:
		where += " AND deleted = $4"

		args = append(args, filter.Deleted)
	case filter.CreatedBetween[0].Before(filter.CreatedBetween[1]):
		where += " AND created_at BETWEEN $5 AND $6"

		args = append(args, filter.CreatedBetween[0], filter.CreatedBetween[1])
	}

	findAllQuery := fmt.Sprintf(`SELECT * FROM %s WHERE %s ORDER BY id`, userTable, where)

	results, err := query[domain.User](ctx, r.db, findAllQuery, args...)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(results))

	for i, r := range results {
		users[i] = *r
	}

	return users, nil
}

// FindByID returns a user from the database by id
func (r *UserRepo) FindByID(ctx context.Context, usrID string) (*domain.User, error) {
	findByIDQuery := fmt.Sprintf(`SELECT * FROM %s
	WHERE id = $1 AND NOT deleted AND active`, userTable)

	user, err := queryRow[domain.User](ctx, r.db, findByIDQuery, usrID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewUserNotFoundError(usrID)
		}

		return nil, fmt.Errorf("find user by ID error: %w", err)
	}

	return user, nil
}

// FindByUsername returns a user from the database by username
func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	findByUsernameQuery := fmt.Sprintf(`SELECT * FROM %s
	WHERE username = $1 AND NOT deleted AND active`, userTable)

	user, err := queryRow[domain.User](ctx, r.db, findByUsernameQuery, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewUserNotFoundError(username)
		}

		return nil, fmt.Errorf("find user by username error: %w", err)
	}

	return user, nil
}

// Create a new user in the database
func (r *UserRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	createUserQuery := fmt.Sprintf(`INSERT INTO %s (
		username, email, password, first_name, last_name
	) VALUES (
		$1, $2, $3, $4, $5
	) RETURNING *`, userTable)

	newUsr, err := queryRow[domain.User](
		ctx,
		r.db,
		createUserQuery,
		user.Username,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
	)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}

	return newUsr, nil
}

// Update a user in the database
func (r *UserRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	updateUserQuery := fmt.Sprintf(`UPDATE %s SET
			username = $1,
			email = $2,
			password = $3,
			first_name = $4,
			last_name = $5
		WHERE id = $6
		RETURNING *`, userTable)

	updatedUsr, err := queryRow[domain.User](
		ctx,
		r.db,
		updateUserQuery,
		user.Username,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("update user error: %w", err)
	}

	return updatedUsr, nil
}

// SoftDelete a user in the database
func (r *UserRepo) SoftDelete(ctx context.Context, usrID string) error {
	softDeleteUserQuery := fmt.Sprintf(`UPDATE %s SET
		deleted = true,
		deleted_at = NOW()
	WHERE id = $1`, userTable)

	_, err := exec(ctx, r.db, softDeleteUserQuery, usrID)
	if err != nil {
		return fmt.Errorf("soft delete user error: %w", err)
	}

	return nil
}

// Delete a user from the database
func (r *UserRepo) Delete(ctx context.Context, id string) error {
	deleteUserQuery := fmt.Sprintf(`DELETE FROM %s
	WHERE id = $1`, userTable)

	_, err := exec(ctx, r.db, deleteUserQuery, id)
	if err != nil {
		return fmt.Errorf("delete user error: %w", err)
	}

	return nil
}

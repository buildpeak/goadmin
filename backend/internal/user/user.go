package user

import (
	"context"
	"fmt"

	"goadmin-backend/internal/domain"
)

type Service interface {
	List(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
}

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService( //nolint: ireturn // it's a factory function
	userRepo domain.UserRepository,
) Service {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) List(
	ctx context.Context,
	filter *domain.UserFilter,
) ([]domain.User, error) {
	users, err := s.userRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find all users error: %w", err)
	}

	return users, nil
}

func (s *userService) GetByID(
	ctx context.Context,
	id string,
) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find user by id error: %w", err)
	}

	return user, nil
}

func (s *userService) Update(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {
	updated, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("update user error: %w", err)
	}

	return updated, nil
}

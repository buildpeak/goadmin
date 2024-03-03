package user

import (
	"context"
	"fmt"

	"goadmin-backend/internal/domain"
)

type Service interface {
	List(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error)
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

package user

import (
	"context"

	"goadmin-backend/internal/domain"
)

type UserService interface {
	List(ctx context.Context, filter domain.UserFilter) ([]domain.User, error)
}

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) List(ctx context.Context, filter domain.UserFilter) ([]domain.User, error) {
	return s.userRepo.FindAll(ctx, filter)
}

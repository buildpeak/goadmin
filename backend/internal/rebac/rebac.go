package rebac

import (
	"context"

	"goadmin-backend/internal/domain"
)

type ReBACCheckResult struct {
	Allowed bool
	Depth   int
}

func Check(
	ctx context.Context,
	user *domain.User,
	resource *domain.Entity,
	action string,
) (*ReBACCheckResult, error) {
	depth := 0
	allowed := false

	return &ReBACCheckResult{
		Allowed: allowed,
		Depth:   depth,
	}, nil
}

type ReBACService interface {
	FindRelationTuple(
		ctx context.Context,
		user *domain.User,
		resource *domain.Entity,
		action string,
	) (*domain.RelationTuple, error)
}

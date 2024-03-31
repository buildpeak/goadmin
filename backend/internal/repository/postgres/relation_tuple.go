package postgres

import (
	"context"
	"fmt"
	"goadmin-backend/internal/domain"
)

type RelationTupleRepo struct {
	db Queryer
}

func NewRelationTupleRepo(db Queryer) *RelationTupleRepo {
	return &RelationTupleRepo{
		db: db,
	}
}

func (rtr *RelationTupleRepo) FindRelationTuple(
	ctx context.Context,
	user *domain.User,
	entity *domain.Entity,
	action string,
) ([]*domain.RelationTuple, error) {
	sql := fmt.Sprintf(`
		SELECT
			*
		FROM
			%s
		WHERE entity_type = $1
			AND entity_id = $2
			AND action = $3
	`, relationTupleTable)

	rt, err := query[domain.RelationTuple](ctx, rtr.db, sql, entity.EntityType, entity.EntityID, action)
	if err != nil {
		return nil, err
	}

	return rt, nil
}

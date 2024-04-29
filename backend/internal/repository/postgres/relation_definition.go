package postgres

import (
	"context"
	"fmt"

	"goadmin-backend/internal/domain"
)

type RelationDefinitionRepo struct {
	db Queryer
}

func NewRelationDefinitionRepo(db Queryer) *RelationDefinitionRepo {
	return &RelationDefinitionRepo{
		db: db,
	}
}

func (rdr *RelationDefinitionRepo) FindRelationDefinition(
	ctx context.Context,
	entityType string,
) ([]*domain.RelationDefinition, error) {
	sql := fmt.Sprintf(`
		SELECT
			*
		FROM
			%s
		WHERE entity_type = $1
	`, relationDefinition)

	rd, err := query[domain.RelationDefinition](ctx, rdr.db, sql, entityType)
	if err != nil {
		return nil, err
	}

	return rd, nil
}

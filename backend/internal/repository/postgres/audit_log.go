package postgres

import (
	"context"
	"fmt"

	"goadmin-backend/internal/domain"
)

var _ domain.AuditLogRepository = &AuditLogRepo{}

type AuditLogRepo struct {
	db Queryer
}

func NewAuditLogRepo(db Queryer) *AuditLogRepo {
	return &AuditLogRepo{db: db}
}

func (r *AuditLogRepo) Create(
	ctx context.Context,
	log *domain.AuditLog,
) (*domain.AuditLog, error) {
	sqlQuery := fmt.Sprintf(
		`INSERT INTO %s (actor_id, actor_username, action, resource_type, resource_id, detail, ip_address)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`,
		auditLogTable,
	)

	result, err := queryRow[domain.AuditLog](
		ctx, r.db, sqlQuery,
		log.ActorID, log.ActorUsername, log.Action,
		log.ResourceType, log.ResourceID, log.Detail, log.IPAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("create audit log error: %w", err)
	}

	return result, nil
}

func (r *AuditLogRepo) FindAll(
	ctx context.Context,
	filter *domain.AuditLogFilter,
) ([]domain.AuditLog, error) {
	if filter == nil {
		filter = &domain.AuditLogFilter{}
	}

	where := "1 = 1"
	args := []interface{}{}

	if filter.ActorUsername != "" {
		args = append(args, filter.ActorUsername)
		where += fmt.Sprintf(" AND actor_username = $%d", len(args))
	}

	if filter.Action != "" {
		args = append(args, filter.Action)
		where += fmt.Sprintf(" AND action = $%d", len(args))
	}

	if filter.ResourceType != "" {
		args = append(args, filter.ResourceType)
		where += fmt.Sprintf(" AND resource_type = $%d", len(args))
	}

	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	args = append(args, filter.Limit, filter.Offset)
	sqlQuery := fmt.Sprintf(
		`SELECT * FROM %s WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		auditLogTable, where, len(args)-1, len(args),
	)

	results, err := query[domain.AuditLog](ctx, r.db, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("find audit logs error: %w", err)
	}

	logs := make([]domain.AuditLog, len(results))
	for i, r := range results {
		logs[i] = *r
	}

	return logs, nil
}

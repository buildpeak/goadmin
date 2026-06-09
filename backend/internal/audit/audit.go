package audit

import (
	"context"
	"fmt"

	"goadmin-backend/internal/domain"
)

type Service interface {
	Log(ctx context.Context, entry *domain.AuditLog) (*domain.AuditLog, error)
	List(ctx context.Context, filter *domain.AuditLogFilter) ([]domain.AuditLog, error)
}

type auditService struct {
	repo domain.AuditLogRepository
}

func NewService(repo domain.AuditLogRepository) Service {
	return &auditService{repo: repo}
}

func (s *auditService) Log(
	ctx context.Context,
	entry *domain.AuditLog,
) (*domain.AuditLog, error) {
	result, err := s.repo.Create(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("audit log error: %w", err)
	}

	return result, nil
}

func (s *auditService) List(
	ctx context.Context,
	filter *domain.AuditLogFilter,
) ([]domain.AuditLog, error) {
	logs, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list audit logs error: %w", err)
	}

	return logs, nil
}

package domain

import (
	"context"
	"time"
)

type AuditLog struct {
	ID            string    `json:"id"`
	ActorID       string    `json:"actor_id,omitempty"`
	ActorUsername  string    `json:"actor_username"`
	Action        string    `json:"action"`
	ResourceType  string    `json:"resource_type"`
	ResourceID    string    `json:"resource_id,omitempty"`
	Detail        string    `json:"detail,omitempty"`
	IPAddress     string    `json:"ip_address,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type AuditLogFilter struct {
	ActorUsername string `json:"actor_username"`
	Action        string `json:"action"`
	ResourceType  string `json:"resource_type"`
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *AuditLog) (*AuditLog, error)
	FindAll(ctx context.Context, filter *AuditLogFilter) ([]AuditLog, error)
}

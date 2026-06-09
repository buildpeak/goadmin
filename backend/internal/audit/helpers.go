package audit

import (
	"context"
	"net"
	"net/http"

	"goadmin-backend/internal/domain"
)

func LogFromRequest(
	ctx context.Context,
	req *http.Request,
	svc Service,
	actorID, actorUsername, action, resourceType, resourceID, detail string,
) {
	if actorUsername == "" {
		actorUsername = "anonymous"
	}

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	if ip == "" {
		ip = req.RemoteAddr
	}

	_, _ = svc.Log(ctx, &domain.AuditLog{
		ActorID:       actorID,
		ActorUsername:  actorUsername,
		Action:        action,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		Detail:        detail,
		IPAddress:     ip,
	})
}

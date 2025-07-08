package auditlog

import (
	"context"

	"github.com/audit-log-service/internal/model"
	"github.com/audit-log-service/internal/repository"
)

// Controller provides specifications related to Audit Log functionality.
type Controller interface {
	// HandleAuditLogMessage processes a Kafka message and saves the audit log entry.
	HandleAuditLogMessage(ctx context.Context, log model.AuditLog) error
	// GetAuditLogs returns paginated audit logs based on filter criteria.
	GetAuditLogs(ctx context.Context, criteria model.FilterCriteria) ([]model.AuditLog, error)
}

// impl holds reference
// to domain repository implementation reference.
type impl struct {
	repo repository.Registry
}

// New returns an implementation instance satisfying Repository
func New(
	repo repository.Registry,
) Controller {
	return impl{
		repo: repo,
	}
}

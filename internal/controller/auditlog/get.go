package auditlog

import (
	"context"

	"github.com/audit-log-service/internal/model"
)

// GetAuditLogs returns paginated audit logs based on filter criteria.
func (i impl) GetAuditLogs(ctx context.Context, criteria model.FilterCriteria) ([]model.AuditLog, error) {
	return i.repo.GetAuditLog().GetAuditLogs(ctx, criteria)
}

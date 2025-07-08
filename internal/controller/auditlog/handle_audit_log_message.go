package auditlog

import (
	"context"

	"github.com/audit-log-service/internal/model"
)

// HandleAuditLogMessage processes a Kafka message and saves the audit log entry.
func (i impl) HandleAuditLogMessage(ctx context.Context, log model.AuditLog) error {
	if _, err := i.repo.GetAuditLog().InsertAuditLog(ctx, log); err != nil {
		return err
	}
	return nil
}

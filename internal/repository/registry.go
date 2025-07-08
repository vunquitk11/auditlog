package repository

import (
	"github.com/audit-log-service/internal/repository/auditlog"
)

// GetAuditLog returns Audit Log repository
func (i impl) GetAuditLog() auditlog.Repository {
	return auditlog.New(i.db)
}

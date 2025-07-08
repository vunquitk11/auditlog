package auditlog

import (
	"context"

	"github.com/audit-log-service/internal/model"
	"gorm.io/gorm"
)

// Repository interface provides the specification for Audit Log feature functionality.
type Repository interface {
	// InsertAuditLog inserts a new audit log entry into the database.
	InsertAuditLog(ctx context.Context, log model.AuditLog) (model.AuditLog, error)
	// GetAuditLogs retrieves paginated audit log records given the criteria.
	GetAuditLogs(ctx context.Context, criteria model.FilterCriteria) ([]model.AuditLog, error)
}

type impl struct {
	db *gorm.DB
}

// New returns an implementation instance satisfying Repository
func New(db *gorm.DB) Repository {
	return &impl{db: db}
}

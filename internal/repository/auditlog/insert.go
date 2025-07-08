package auditlog

import (
	"context"

	"github.com/audit-log-service/internal/model"
	pkgerrors "github.com/pkg/errors"
)

// InsertAuditLog inserts a new audit log entry into the database.
func (i *impl) InsertAuditLog(ctx context.Context, log model.AuditLog) (model.AuditLog, error) {
	if err := i.db.WithContext(ctx).Create(&log).Error; err != nil {
		return model.AuditLog{}, pkgerrors.WithStack(err)
	}
	return log, nil
}

package auditlog

import (
	"context"

	"github.com/audit-log-service/internal/model"
	pkgerrors "github.com/pkg/errors"
)

// GetAuditLogs retrieves paginated audit log records given the criteria.
func (i *impl) GetAuditLogs(ctx context.Context, criteria model.FilterCriteria) ([]model.AuditLog, error) {
	db := i.db.WithContext(ctx)
	if criteria.ServiceName != "" {
		db = db.Where("service_name = ?", criteria.ServiceName)
	}
	if criteria.Username != "" {
		db = db.Where("username = ?", criteria.Username)
	}
	if criteria.Action != "" {
		db = db.Where("action = ?", criteria.Action)
	}
	if criteria.ResourceType != "" {
		db = db.Where("resource_type = ?", criteria.ResourceType)
	}
	if criteria.ResourceID != "" {
		db = db.Where("resource_id = ?", criteria.ResourceID)
	}
	if !criteria.TimestampFrom.IsZero() {
		db = db.Where("timestamp >= ?", criteria.TimestampFrom)
	}
	if !criteria.TimestampTo.IsZero() {
		db = db.Where("timestamp <= ?", criteria.TimestampTo)
	}
	db = db.Order("timestamp DESC").Limit(criteria.GetLimit()).Offset(criteria.GetOffset())
	var logs []model.AuditLog
	if err := db.Find(&logs).Error; err != nil {
		return nil, pkgerrors.WithStack(err)
	}
	return logs, nil
}

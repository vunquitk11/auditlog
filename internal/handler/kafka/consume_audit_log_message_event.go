package kafka

import (
	"context"
	"log"

	"github.com/audit-log-service/internal/model"
	"github.com/audit-log-service/pkg/metrics"
)

// ConsumeAuditLogMessageEvent processes a single audit log message (already parsed from JSON).
func (h *Handler) ConsumeAuditLogMessageEvent(ctx context.Context, auditMsg *model.AuditMessage) error {
	auditLog := auditMsg.ToAuditLog()
	err := h.auditLogCtrl.HandleAuditLogMessage(ctx, *auditLog)
	if err != nil {
		metrics.AuditLogsFailedTotal.WithLabelValues(auditMsg.ServiceName, auditMsg.Action).Inc()
		log.Printf("Failed to process audit log: %v", err)
		return err
	}
	metrics.AuditLogsTotal.WithLabelValues(auditMsg.ServiceName, auditMsg.Action).Inc()
	log.Printf("Audit log processed: id=%s", auditLog.ID)
	return nil
}

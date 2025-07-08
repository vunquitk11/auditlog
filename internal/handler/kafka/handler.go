package kafka

import (
	"github.com/audit-log-service/internal/controller/auditlog"
)

// Handler is responsible for processing audit log messages.
type Handler struct {
	auditLogCtrl auditlog.Controller
}

// New returns a new instance of Handler.
func New(auditLogCtrl auditlog.Controller) Handler {
	return Handler{
		auditLogCtrl: auditLogCtrl,
	}
}

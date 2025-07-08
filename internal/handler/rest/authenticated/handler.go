package authenticated

import (
	"github.com/audit-log-service/internal/controller/auditlog"
)

// Handler is the web handler for this pkg
type Handler struct {
	auditLogCtrl auditlog.Controller
}

// New instantiates a new Handler and returns it
func New(
	auditLogCtrl auditlog.Controller,
) Handler {
	return Handler{
		auditLogCtrl: auditLogCtrl,
	}
}

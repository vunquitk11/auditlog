package repository

import (
	"gorm.io/gorm"

	"github.com/audit-log-service/internal/repository/auditlog"
)

type Registry interface {
	GetAuditLog() auditlog.Repository
}

type impl struct {
	db           *gorm.DB
	auditLogRepo auditlog.Repository
}

func New(db *gorm.DB) Registry {
	return &impl{
		db: db,
	}
}

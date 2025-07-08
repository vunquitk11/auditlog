package model

import (
	"time"

	"gorm.io/datatypes"
)

// AuditLog represents an audit log entry in the system.
type AuditLog struct {
	ID           string         `gorm:"primaryKey"`
	ServiceName  string         `gorm:"size:255;not null"`
	Username     string         `gorm:"size:255;not null"`
	Action       string         `gorm:"size:255;not null"`
	ResourceType string         `gorm:"size:255;not null"`
	ResourceID   string         `gorm:"size:255;not null"`
	Timestamp    time.Time      `gorm:"not null"`
	BeforeState  datatypes.JSON `gorm:"type:jsonb"`
	AfterState   datatypes.JSON `gorm:"type:jsonb"`
	Metadata     datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
}

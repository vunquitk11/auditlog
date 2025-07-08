package model

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// AuditMessage represents the structure of audit log messages sent to Kafka
type AuditMessage struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Action      string                 `json:"action"`
	Resource    string                 `json:"resource"`
	ResourceID  string                 `json:"resource_id"`
	Details     map[string]interface{} `json:"details"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Timestamp   time.Time              `json:"timestamp"`
	ServiceName string                 `json:"service_name"`
}

// ToJSON converts the audit message to JSON bytes
func (am *AuditMessage) ToJSON() ([]byte, error) {
	return json.Marshal(am)
}

// FromJSON creates an audit message from JSON bytes
func FromJSON(data []byte) (*AuditMessage, error) {
	var msg AuditMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// ToAuditLog converts the audit message to an AuditLog model
func (am *AuditMessage) ToAuditLog() *AuditLog {
	// Convert details to JSON
	detailsJSON, _ := json.Marshal(am.Details)

	// Create metadata with additional info
	metadata := map[string]interface{}{
		"ip_address": am.IPAddress,
		"user_agent": am.UserAgent,
	}
	metadataJSON, _ := json.Marshal(metadata)

	return &AuditLog{
		ID:           am.ID,
		ServiceName:  am.ServiceName,
		Username:     am.UserID,
		Action:       am.Action,
		ResourceType: am.Resource,
		ResourceID:   am.ResourceID,
		Timestamp:    am.Timestamp,
		BeforeState:  datatypes.JSON("{}"),
		AfterState:   detailsJSON,
		Metadata:     metadataJSON,
	}
}

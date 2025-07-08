package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/audit-log-service/internal/model"
)

// Producer handles sending audit log messages to Kafka
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 10 * time.Second

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

// SendAuditMessage sends an audit message to Kafka
func (p *Producer) SendAuditMessage(message *model.AuditMessage) error {
	jsonData, err := message.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       sarama.StringEncoder(message.ID),
		Value:     sarama.ByteEncoder(jsonData),
		Timestamp: message.Timestamp,
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d", p.topic, partition, offset)
	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.producer.Close()
}

// SendAuditLog is a convenience method to send audit log data directly
func (p *Producer) SendAuditLog(auditLog *model.AuditLog) error {
	// Convert AuditLog to AuditMessage
	var details map[string]interface{}
	if len(auditLog.AfterState) > 0 {
		_ = json.Unmarshal(auditLog.AfterState, &details)
	}
	var metadata map[string]interface{}
	if len(auditLog.Metadata) > 0 {
		_ = json.Unmarshal(auditLog.Metadata, &metadata)
	}
	message := &model.AuditMessage{
		ID:          auditLog.ID,
		UserID:      auditLog.Username,
		Action:      auditLog.Action,
		Resource:    auditLog.ResourceType,
		ResourceID:  auditLog.ResourceID,
		Details:     details,
		IPAddress:   fmt.Sprintf("%v", metadata["ip_address"]),
		UserAgent:   fmt.Sprintf("%v", metadata["user_agent"]),
		Timestamp:   auditLog.Timestamp,
		ServiceName: auditLog.ServiceName,
	}
	return p.SendAuditMessage(message)
}

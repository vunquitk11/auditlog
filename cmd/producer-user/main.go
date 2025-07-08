package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/audit-log-service/internal/model"
	"github.com/audit-log-service/pkg/kafka"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const (
	topicName = "audit-logs"
	userAgent = "UserProducer/1.0"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Parse command line flags
	var (
		brokers  = flag.String("brokers", getEnv("KAFKA_BROKERS", "localhost:9092"), "Kafka brokers (comma-separated)")
		interval = flag.Duration("interval", 5*time.Second, "Interval between messages")
		count    = flag.Int("count", 0, "Number of messages to send (0 for infinite)")
	)
	flag.Parse()

	// Parse brokers
	brokerList := strings.Split(*brokers, ",")
	for i, broker := range brokerList {
		brokerList[i] = strings.TrimSpace(broker)
	}

	// Create Kafka producer
	producer, err := kafka.NewProducer(brokerList, topicName)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	log.Printf("User producer started. Sending messages to topic: %s", topicName)
	log.Printf("Brokers: %v", brokerList)
	log.Printf("Interval: %v", *interval)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutdown signal received, stopping producer...")
		cancel()
	}()

	// Start producing messages
	messageCount := 0
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Producer stopped. Total messages sent: %d", messageCount)
			return
		case <-ticker.C:
			if err := sendUserAuditMessage(producer); err != nil {
				log.Printf("Error sending message: %v", err)
				continue
			}
			messageCount++

			if *count > 0 && messageCount >= *count {
				log.Printf("Reached message count limit (%d). Stopping producer.", *count)
				return
			}
		}
	}
}

// sendUserAuditMessage sends a simulated user action audit message
func sendUserAuditMessage(producer *kafka.Producer) error {
	// Simulate different user actions
	actions := []string{"login", "logout", "create_user", "update_profile", "delete_account", "change_password"}
	resources := []string{"user", "profile", "account", "password"}

	action := actions[rand.Intn(len(actions))]
	resource := resources[rand.Intn(len(resources))]

	// Create audit message
	message := &model.AuditMessage{
		ID:         uuid.New().String(),
		UserID:     fmt.Sprintf("user_%d", rand.Intn(1000)+1),
		Action:     action,
		Resource:   resource,
		ResourceID: fmt.Sprintf("res_%d", rand.Intn(10000)+1),
		Details: map[string]interface{}{
			"ip_address": fmt.Sprintf("192.168.1.%d", rand.Intn(255)+1),
			"session_id": uuid.New().String(),
			"success":    rand.Float32() > 0.1, // 90% success rate
		},
		IPAddress:   fmt.Sprintf("192.168.1.%d", rand.Intn(255)+1),
		UserAgent:   userAgent,
		Timestamp:   time.Now(),
		ServiceName: "user-service",
	}

	// Send message
	if err := producer.SendAuditMessage(message); err != nil {
		return fmt.Errorf("failed to send user audit message: %w", err)
	}

	log.Printf("Sent user audit message: %s - %s %s", message.UserID, message.Action, message.Resource)
	return nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

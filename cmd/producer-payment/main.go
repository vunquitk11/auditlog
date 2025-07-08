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
	userAgent = "PaymentProducer/1.0"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Parse command line flags
	var (
		brokers  = flag.String("brokers", getEnv("KAFKA_BROKERS", "localhost:9092"), "Kafka brokers (comma-separated)")
		interval = flag.Duration("interval", 3*time.Second, "Interval between messages")
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

	log.Printf("Payment producer started. Sending messages to topic: %s", topicName)
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
			if err := sendPaymentAuditMessage(producer); err != nil {
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

// sendPaymentAuditMessage sends a simulated payment action audit message
func sendPaymentAuditMessage(producer *kafka.Producer) error {
	// Simulate different payment actions
	actions := []string{"payment_initiated", "payment_processed", "payment_failed", "refund_issued", "chargeback_received"}
	resources := []string{"payment", "transaction", "order", "refund", "chargeback"}

	action := actions[rand.Intn(len(actions))]
	resource := resources[rand.Intn(len(resources))]

	// Generate random amount
	amount := float64(rand.Intn(10000)+1) / 100.0 // $0.01 to $100.00

	// Create audit message
	message := &model.AuditMessage{
		ID:         uuid.New().String(),
		UserID:     fmt.Sprintf("user_%d", rand.Intn(1000)+1),
		Action:     action,
		Resource:   resource,
		ResourceID: fmt.Sprintf("txn_%d", rand.Intn(100000)+1),
		Details: map[string]interface{}{
			"amount":         amount,
			"currency":       "USD",
			"payment_method": getRandomPaymentMethod(),
			"merchant_id":    fmt.Sprintf("merchant_%d", rand.Intn(100)+1),
			"order_id":       fmt.Sprintf("order_%d", rand.Intn(1000000)+1),
			"status":         getRandomStatus(),
		},
		IPAddress:   fmt.Sprintf("10.0.0.%d", rand.Intn(255)+1),
		UserAgent:   userAgent,
		Timestamp:   time.Now(),
		ServiceName: "payment-service",
	}

	// Send message
	if err := producer.SendAuditMessage(message); err != nil {
		return fmt.Errorf("failed to send payment audit message: %w", err)
	}

	log.Printf("Sent payment audit message: %s - %s %s (%.2f USD)",
		message.UserID, message.Action, message.Resource, amount)
	return nil
}

// getRandomPaymentMethod returns a random payment method
func getRandomPaymentMethod() string {
	methods := []string{"credit_card", "debit_card", "paypal", "bank_transfer", "crypto"}
	return methods[rand.Intn(len(methods))]
}

// getRandomStatus returns a random payment status
func getRandomStatus() string {
	statuses := []string{"pending", "completed", "failed", "cancelled", "refunded"}
	return statuses[rand.Intn(len(statuses))]
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

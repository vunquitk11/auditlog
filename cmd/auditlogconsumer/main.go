package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/audit-log-service/cmd/common"
	"github.com/audit-log-service/internal/controller/auditlog"
	kafkaHandler "github.com/audit-log-service/internal/handler/kafka"
	"github.com/audit-log-service/internal/model"
	"github.com/audit-log-service/pkg/kafka"
	"github.com/audit-log-service/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Register custom metrics with Prometheus
	metrics.RegisterMetrics()
	// Create dummy series to display the metric even when there is no real data yet
	metrics.AuditLogsTotal.WithLabelValues("startup", "startup").Add(0)

	// Init repository connection
	repoRegistry := common.InitRepository()

	// Init audit log controller
	auditCtrl := auditlog.New(repoRegistry)
	handler := kafkaHandler.New(auditCtrl)

	// Expose /metrics endpoint for Prometheus
	go func() {
		log.Println("Starting metrics endpoint at :2112/metrics")
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "audit-logs"
	}
	groupID := os.Getenv("KAFKA_CONSUMER_GROUP")
	if groupID == "" {
		groupID = "audit-log-consumer"
	}

	consumer, err := kafka.NewConsumerGroup(brokers, groupID)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer group: %v", err)
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	cgHandler := &auditLogConsumerGroupHandler{handler: &handler}

	go func() {
		<-sigchan
		log.Println("Shutting down Kafka consumer...")
		cancel()
	}()

	for {
		if err := consumer.Consume(ctx, []string{topic}, cgHandler); err != nil {
			log.Printf("Error from consumer: %v", err)
			break
		}
		if ctx.Err() != nil {
			break
		}
	}

	log.Println("Audit log consumer stopped.")
}

// auditLogConsumerGroupHandler implements sarama.ConsumerGroupHandler
type auditLogConsumerGroupHandler struct {
	handler *kafkaHandler.Handler
}

func (h *auditLogConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *auditLogConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *auditLogConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Received message: topic=%s partition=%d offset=%d", message.Topic, message.Partition, message.Offset)
		auditMsg, err := model.FromJSON(message.Value)
		if err != nil {
			log.Printf("Failed to parse audit message: %v", err)
			continue
		}
		err = h.handler.ConsumeAuditLogMessageEvent(session.Context(), auditMsg)
		if err != nil {
			log.Printf("Failed to process audit log: %v", err)
			continue
		}
		session.MarkMessage(message, "")
	}
	return nil
}

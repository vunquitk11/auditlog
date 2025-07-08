package kafka

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

// NewConsumerGroup creates a new Kafka consumer group with default balanced configuration.
func NewConsumerGroup(brokers string, groupID string) (sarama.ConsumerGroup, error) {
	brokerList := strings.Split(brokers, ",")
	for i, broker := range brokerList {
		brokerList[i] = strings.TrimSpace(broker)
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	consumer, err := sarama.NewConsumerGroup(brokerList, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return consumer, nil
}

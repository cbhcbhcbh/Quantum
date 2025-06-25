package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
)

type ConsumerGroup struct {
	sarama.ConsumerGroup
	groupID string
	topics  []string
	Ready   chan bool
}

func NewConsumerGroup(ctx context.Context, topics, addrs []string, groupID string) *ConsumerGroup {
	cg := &ConsumerGroup{
		groupID: groupID,
		topics:  topics,
		Ready:   make(chan bool),
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		log.C(ctx).Errorw("Failed to start Sarama consumer group", "error", err)
		return nil
	}
	cg.ConsumerGroup = consumerGroup

	return cg
}

func (cg *ConsumerGroup) Consume(ctx context.Context, handler sarama.ConsumerGroupHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := cg.ConsumerGroup.Consume(ctx, cg.topics, handler); err != nil {
				log.C(ctx).Errorw("Error from consumer group", "error", err)
				return err
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}
	}
}

func (cg *ConsumerGroup) Close(ctx context.Context) {
	err := cg.ConsumerGroup.Close()
	if err != nil {
		log.C(ctx).Errorw("Failed to close Sarama consumergroup", "error", err)
	}

	close(cg.Ready)
}

package kafka

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
)

type KafkaConsumer struct {
	Address []string
	Topic   []string
	GroupID string
}

type Consumer struct {
	addr          []string
	WG            sync.WaitGroup
	PartitionList []int32
	Consumer      sarama.Consumer
}

func NewConsumer(ctx context.Context, addr []string, topic string) *Consumer {
	c := Consumer{}
	c.addr = addr

	consumer, err := sarama.NewConsumer(c.addr, nil)
	if err != nil {
		log.C(ctx).Errorw("Failed to start Sarama producer", "error", err)
		return nil
	}
	c.Consumer = consumer

	partitionList, err := c.Consumer.Partitions(topic)
	if err != nil {
		log.C(ctx).Errorw("Failed to get partitions for topic", "topic", topic, "error", err)
		return nil
	}

	c.PartitionList = partitionList

	return &c
}

func (c *Consumer) Consume(ctx context.Context, topic string, partition int32, handler func(msg *sarama.ConsumerMessage)) {
    pc, err := c.Consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
    if err != nil {
        log.C(ctx).Errorw("Failed to start partition consumer", "topic", topic, "partition", partition, "error", err)
        return
    }
    defer pc.Close()

    for {
        select {
        case msg := <-pc.Messages():
            if msg != nil {
                handler(msg)
            }
        case err := <-pc.Errors():
            if err != nil {
                log.C(ctx).Errorw("Partition consumer error", "topic", topic, "partition", partition, "error", err)
            }
        case <-ctx.Done():
            return
        }
    }
}

func (c *Consumer) Close(ctx context.Context) {
	err := c.Consumer.Close()
	if err != nil {
		log.C(ctx).Errorw("Failed to close Sarama consumer", "error", err)
	}
}

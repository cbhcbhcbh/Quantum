package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
)

var (
	P *Producer
)

type KafkaProducer struct {
	Affress []string
	Topic   string
}

type Producer struct {
	addr     []string
	config   *sarama.Config
	producer sarama.SyncProducer
}

func NewProducer(addr []string) *Producer {
	p := Producer{}
	p.config = sarama.NewConfig()
	p.config.Producer.Return.Successes = true
	p.config.Producer.RequiredAcks = sarama.WaitForAll
	p.config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	p.addr = addr

	producer, err := sarama.NewSyncProducer(p.addr, p.config)
	if err != nil {
		log.C(context.Background()).Errorw("Failed to start Sarama producer", "error", err)
		return nil
	}
	p.producer = producer

	return &p
}

func (p *Producer) Push(buf []byte, key string, topic string) (int32, int64, error) {
	msg := &sarama.ProducerMessage{}
	msg.Key = sarama.StringEncoder(key)
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(buf)
	return p.producer.SendMessage(msg)
}

func (p *Producer) Close(ctx context.Context) {
	err := p.producer.Close()
	if err != nil {
		log.C(ctx).Errorw("Failed to close Sarama producer", "error", err)
	}
}

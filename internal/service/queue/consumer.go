package queue

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/date"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/message"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"github.com/cbhcbhcbh/Quantum/pkg/kafka"
)

type PrivateHandle struct{}

type GroupHandle struct{}

func (pH *PrivateHandle) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) {
	// TODO: Batch multiple messages and insert them into the database in a single operation for better performance.
	log.C(ctx).Infow("Consumed private offline message", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)

	var wsMsg message.WsMessage
	if err := json.Unmarshal(msg.Value, &wsMsg); err != nil {
		log.C(ctx).Errorw("Failed to unmarshal private offline message", "error", err)
		return
	}

	receiveId := wsMsg.ToID

	store.S.OfflineMessage().Create(ctx, &model.OfflineMessageM{
		Status:    0,
		SendTime:  date.TimeUnix(),
		ReceiveID: receiveId,
		Message:   msg.Value,
	})
}

func (gH *GroupHandle) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) {
	// TODO: Batch multiple messages and insert them into the database in a single operation for better performance.
	log.C(ctx).Infow("Consumed group offline message", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)
	var wsMsg message.WsMessage
	if err := json.Unmarshal(msg.Value, &wsMsg); err != nil {
		log.C(ctx).Errorw("Failed to unmarshal group offline message", "error", err)
		return
	}

	receiveId := wsMsg.ToID

	store.S.GroupOfflineMessage().Create(ctx, &model.GroupOfflineMessageM{
		Status:    0,
		SendTime:  date.TimeUnix(),
		ReceiveID: receiveId,
		Message:   msg.Value,
	})
}

func ConsumerMessage(ctx context.Context, addr []string) {
	go func() {
		consumer := kafka.NewConsumer(ctx, addr, known.OfflinePrivateTopic)
		handler := &PrivateHandle{}
		if consumer == nil {
			log.C(ctx).Errorw("Failed to create consumer for private topic")
			return
		}
		for _, partition := range consumer.PartitionList {
			go consumer.Consume(ctx, known.OfflinePrivateTopic, partition, handler.HandleMessage)
		}
	}()

	go func() {
		consumer := kafka.NewConsumer(ctx, addr, known.OfflineGroupTopic)
		handler := &GroupHandle{}
		if consumer == nil {
			log.C(ctx).Errorw("Failed to create consumer for group topic")
			return
		}
		for _, partition := range consumer.PartitionList {
			go consumer.Consume(ctx, known.OfflineGroupTopic, partition, handler.HandleMessage)
		}
	}()
}

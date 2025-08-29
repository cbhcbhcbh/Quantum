package forwarder

import (
	"context"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
	"github.com/cbhcbhcbh/Quantum/pkg/common/util"
	"github.com/cbhcbhcbh/Quantum/pkg/infra"
)

var (
	forwardPrefix = "rc:forward"
)

type Subscribers map[string]struct{}

type ForwardRepo interface {
	RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error
	RemoveChannelSession(ctx context.Context, channelID, userID uint64) error
	GetSubscribers(ctx context.Context, channelID uint64) (Subscribers, error)
	ForwardMessage(ctx context.Context, msg *domain.Message, subscribers Subscribers) error
}

type ForwardRepoImpl struct {
	r infra.RedisCache
	p message.Publisher
}

func NewForwardRepoImpl(r infra.RedisCache, p message.Publisher) *ForwardRepoImpl {
	return &ForwardRepoImpl{r, p}
}

func (repo *ForwardRepoImpl) RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error {
	key := util.ConstructKey(forwardPrefix, channelID)
	return repo.r.HSet(ctx, key, strconv.FormatUint(userID, 10), subscriber)
}

func (repo *ForwardRepoImpl) RemoveChannelSession(ctx context.Context, channelID, userID uint64) error {
	key := util.ConstructKey(forwardPrefix, channelID)
	return repo.r.HDel(ctx, key, strconv.FormatUint(userID, 10))
}

func (repo *ForwardRepoImpl) GetSubscribers(ctx context.Context, channelID uint64) (Subscribers, error) {
	key := util.ConstructKey(forwardPrefix, channelID)
	sessionMap, err := repo.r.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	subscribers := make(Subscribers)
	for _, subscriber := range sessionMap {
		subscribers[subscriber] = struct{}{}
	}
	return subscribers, nil
}

func (repo *ForwardRepoImpl) ForwardMessage(ctx context.Context, msg *domain.Message, subscribers Subscribers) error {
	var err error
	for subscriber := range subscribers {
		err = repo.p.Publish(subscriber, message.NewMessage(
			watermill.NewUUID(),
			msg.Encode(),
		))
		if err != nil {
			return err
		}
	}
	return nil
}
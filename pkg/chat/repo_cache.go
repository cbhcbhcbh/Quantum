package chat

import (
	"context"
	"strconv"

	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
	"github.com/cbhcbhcbh/Quantum/pkg/common/util"
	"github.com/cbhcbhcbh/Quantum/pkg/infra"
)

type UserRepoCache interface {
	AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error
	GetUserByID(ctx context.Context, userID uint64) (*domain.User, error)
	IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error)
	GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
	AddOnlineUser(ctx context.Context, channelID uint64, userID uint64) error
	DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error
	GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
}

type MessageRepoCache interface {
	InsertMessage(ctx context.Context, msg *domain.Message) error
	MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error
	PublishMessage(ctx context.Context, msg *domain.Message) error
	ListMessages(ctx context.Context, channelID uint64, pageStateStr string) ([]*domain.Message, string, error)
}

type ChannelRepoCache interface {
	CreateChannel(ctx context.Context, channelID uint64) (*domain.Channel, error)
	DeleteChannel(ctx context.Context, channelID uint64) error
}

type UserRepoCacheImpl struct {
	r        infra.RedisCache
	userRepo UserRepo
}

func NewUserRepoCacheImpl(r infra.RedisCache, userRepo UserRepo) *UserRepoCacheImpl {
	return &UserRepoCacheImpl{r, userRepo}
}

func (cache *UserRepoCacheImpl) AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error {
	if err := cache.userRepo.AddUserToChannel(ctx, channelID, userID); err != nil {
		return nil
	}
	key := util.ConstructKey(known.ChannelUsersPrefix, channelID)
	return cache.r.HSet(ctx, key, strconv.FormatUint(userID, 10), 1)
}

func (cache *UserRepoCacheImpl) GetUserByID(ctx context.Context, userID uint64) (*domain.User, error) {
	return cache.userRepo.GetUserByID(ctx, userID)
}

func (cache *UserRepoCacheImpl) IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error) {
	key := util.ConstructKey(known.ChannelUsersPrefix, channelID)
	var dummy int
	var err error
	channelExists, userExists, err := cache.r.HGetIfKeyExists(ctx, key, strconv.FormatUint(userID, 10), &dummy)
	if err != nil {
		return false, err
	}
	if channelExists {
		if !userExists {
			return false, nil
		}
		return true, nil
	}

	channelUserIDs, err := cache.userRepo.GetChannelUserIDs(ctx, channelID)
	if err != nil {
		return false, err
	}
	channelUserExist := false
	var args []any
	for _, channelUserID := range channelUserIDs {
		if userID == channelUserID {
			channelUserExist = true
		}
		args = append(args, channelUserID, 1)
	}
	if err := cache.r.HSet(ctx, key, args...); err != nil {
		return channelUserExist, err
	}
	return channelUserExist, nil
}

func (cache *UserRepoCacheImpl) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	key := util.ConstructKey(known.ChannelUsersPrefix, channelID)
	userMap, err := cache.r.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	var userIDs []uint64
	if len(userMap) > 0 {
		for userIDStr := range userMap {
			userID, err := strconv.ParseUint(userIDStr, 10, 64)
			if err != nil {
				return nil, err
			}
			userIDs = append(userIDs, userID)
		}
	}

	userIDs, err = cache.userRepo.GetChannelUserIDs(ctx, channelID)
	if err != nil {
		return nil, err
	}
	var args []any
	for _, userID := range userIDs {
		args = append(args, userID, 1)
	}
	if err := cache.r.HSet(ctx, key, args...); err != nil {
		return userIDs, err
	}
	return userIDs, nil
}

func (cache *UserRepoCacheImpl) AddOnlineUser(ctx context.Context, channelID uint64, userID uint64) error {
	key := util.ConstructKey(known.OnlineUsersPrefix, channelID)
	return cache.r.HSet(ctx, key, strconv.FormatUint(userID, 10), 1)
}

func (cache *UserRepoCacheImpl) DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error {
	key := util.ConstructKey(known.OnlineUsersPrefix, channelID)
	userKey := strconv.FormatUint(userID, 10)
	return cache.r.HDel(ctx, key, userKey)
}

func (cache *UserRepoCacheImpl) GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	key := util.ConstructKey(known.OnlineUsersPrefix, channelID)
	userMap, err := cache.r.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	var userIDs []uint64
	for userIDStr := range userMap {
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}

type MessageRepoCacheImpl struct {
	messageRepo MessageRepo
}

func (cache *MessageRepoCacheImpl) InsertMessage(ctx context.Context, msg *domain.Message) error {
	return cache.messageRepo.InsertMessage(ctx, msg)
}

func (cache *MessageRepoCacheImpl) MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error {
	return cache.messageRepo.MarkMessageSeen(ctx, channelID, messageID)
}

func (cache *MessageRepoCacheImpl) PublishMessage(ctx context.Context, msg *domain.Message) error {
	return cache.messageRepo.PublishMessage(ctx, msg)
}

func (cache *MessageRepoCacheImpl) ListMessages(ctx context.Context, channelID uint64, pageStateStr string) ([]*domain.Message, string, error) {
	return cache.messageRepo.ListMessages(ctx, channelID, pageStateStr)
}

type ChannelRepoCacheImpl struct {
	r           infra.RedisCache
	channelRepo ChannelRepo
}

func NewChannelRepoCacheImpl(r infra.RedisCache, channelRepo ChannelRepo) *ChannelRepoCacheImpl {
	return &ChannelRepoCacheImpl{r, channelRepo}
}

func (cache *ChannelRepoCacheImpl) CreateChannel(ctx context.Context, channelID uint64) (*domain.Channel, error) {
	return cache.channelRepo.CreateChannel(ctx, channelID)
}

func (cache *ChannelRepoCacheImpl) DeleteChannel(ctx context.Context, channelID uint64) error {
	if err := cache.channelRepo.DeleteChannel(ctx, channelID); err != nil {
		return err
	}
	cmds := []infra.RedisCmd{
		{
			OpType: infra.DELETE,
			Payload: infra.RedisDeletePayload{
				Key: util.ConstructKey(known.OnlineUsersPrefix, channelID),
			},
		},
		{
			OpType: infra.DELETE,
			Payload: infra.RedisDeletePayload{
				Key: util.ConstructKey(known.ChannelUsersPrefix, channelID),
			},
		},
	}
	return cache.r.ExecPipeLine(ctx, &cmds)
}

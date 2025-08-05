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
}

func (cache *UserRepoCacheImpl) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
}

func (cache *UserRepoCacheImpl) AddOnlineUser(ctx context.Context, channelID uint64, userID uint64) error {
}

func (cache *UserRepoCacheImpl) DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error {
}

func (cache *UserRepoCacheImpl) GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
}

type ChannelRepoCacheImpl struct {
	r           infra.RedisCache
	channelRepo ChannelRepo
}

func NewChannelRepoCacheImpl(r infra.RedisCache, channelRepo ChannelRepo) *ChannelRepoCacheImpl {
	return &ChannelRepoCacheImpl{r, channelRepo}
}

func (cache *ChannelRepoCacheImpl) CreateChannel(ctx context.Context, channelID uint64) (*domain.Channel, error) {
}

func (cache *ChannelRepoCacheImpl) DeleteChannel(ctx context.Context, channelID uint64) error {
}

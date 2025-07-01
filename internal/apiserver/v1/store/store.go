package store

import (
	"sync"

	"gorm.io/gorm"
)

var (
	once sync.Once
	S    *datastore
)

type IStore interface {
	DB() *gorm.DB
	Users() UsersStore
	GroupMessage() GroupMessageStore
	GroupOfflineMessage() GroupOfflineMessageStore
	GroupUserMessage() GroupUserMessageStore
	Message() MessageStore
	OfflineMessage() OfflineMessageStore
	Group() GroupStore
	GroupUser() GroupUserStore
	Friends() FriendsStore
}

type datastore struct {
	db *gorm.DB
}

var _ IStore = (*datastore)(nil)

func NewStore(db *gorm.DB) *datastore {
	once.Do(func() {
		S = &datastore{
			db: db,
		}
	})

	return S
}

func (ds *datastore) DB() *gorm.DB {
	return ds.db
}

func (ds *datastore) Users() UsersStore {
	return NewUsers(ds.db)
}

func (ds *datastore) GroupMessage() GroupMessageStore {
	return NewGroupMessage(ds.db)
}

func (ds *datastore) GroupOfflineMessage() GroupOfflineMessageStore {
	return NewGroupOfflineMessage(ds.db)
}

func (ds *datastore) GroupUserMessage() GroupUserMessageStore {
	return NewGroupUserMessage(ds.db)
}

func (ds *datastore) Message() MessageStore {
	return NewMessage(ds.db)
}

func (ds *datastore) OfflineMessage() OfflineMessageStore {
	return NewOfflineMessage(ds.db)
}

func (ds *datastore) Group() GroupStore {
	return NewGroup(ds.db)
}

func (ds *datastore) GroupUser() GroupUserStore {
	return NewGroupUser(ds.db)
}

func (ds *datastore) Friends() FriendsStore {
	return NewFriends(ds.db)
}

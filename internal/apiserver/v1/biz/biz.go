package biz

import (
	friendrecord "github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz/friendRecord"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz/friends"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz/group"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz/message"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz/session"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz/users"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
)

type IBiz interface {
	Users() users.UserBiz
	Friends() friends.FriendBiz
	Groups() group.GroupBiz
	Sessions() session.SessionBiz
	FriendRecord() friendrecord.FriendRecordBiz
	Message() message.MessageBiz
}

type biz struct {
	ds store.IStore
}

var _ IBiz = (*biz)(nil)

func NewBiz(ds store.IStore) IBiz {
	return &biz{
		ds: ds,
	}
}

func (b *biz) Users() users.UserBiz {
	return users.New(b.ds)
}

func (b *biz) Friends() friends.FriendBiz {
	return friends.New(b.ds)
}

func (b *biz) Groups() group.GroupBiz {
	return group.New(b.ds)
}

func (b *biz) Sessions() session.SessionBiz {
	return session.New(b.ds)
}

func (b *biz) FriendRecord() friendrecord.FriendRecordBiz {
	return friendrecord.New(b.ds)
}

func (b *biz) Message() message.MessageBiz {
	return message.New(b.ds)
}

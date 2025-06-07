package biz

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/biz/users"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/store"
)

type IBiz interface {
	Users() users.UserBiz
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

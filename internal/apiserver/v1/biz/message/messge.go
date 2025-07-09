package message

import "github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"

type MessageBiz interface {
}

type messageBiz struct {
	ds store.IStore
}

var _ MessageBiz = (*messageBiz)(nil)

func New(ds store.IStore) MessageBiz {
	return &messageBiz{ds: ds}
}

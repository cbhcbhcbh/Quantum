package message

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
)

type MessageController struct {
	b biz.IBiz
}

func New(ds store.IStore) *MessageController {
	return &MessageController{
		b: biz.NewBiz(ds),
	}
}

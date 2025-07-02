package session

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
)

type SessionController struct {
	b biz.IBiz
}

func New(ds store.IStore) *SessionController {
	return &SessionController{
		b: biz.NewBiz(ds),
	}
}

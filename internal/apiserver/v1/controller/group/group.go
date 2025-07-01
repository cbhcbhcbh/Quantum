package group

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
)

type GroupController struct {
	b biz.IBiz
}

func New(ds store.IStore) *GroupController {
	return &GroupController{
		b: biz.NewBiz(ds),
	}
}

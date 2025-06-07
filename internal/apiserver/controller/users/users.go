package users

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/biz"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/store"
)

type UserController struct {
	b biz.IBiz
}

func New(ds store.IStore) *UserController {
	return &UserController{
		b: biz.NewBiz(ds),
	}
}

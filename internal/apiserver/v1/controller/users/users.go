package users

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
)

type UserController struct {
	b biz.IBiz
}

func New(ds store.IStore) *UserController {
	return &UserController{
		b: biz.NewBiz(ds),
	}
}

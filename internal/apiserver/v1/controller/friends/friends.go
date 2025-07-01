package friends

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/biz"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
)

type FriendController struct {
	b biz.IBiz
}

func New(ds store.IStore) *FriendController {
	return &FriendController{
		b: biz.NewBiz(ds),
	}
}

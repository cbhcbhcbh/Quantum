package handler

import (
	"github.com/cbhcbhcbh/Quantum/internal/pkg/code"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	"github.com/cbhcbhcbh/Quantum/internal/service/client"
	"github.com/cbhcbhcbh/Quantum/pkg/ws"
	"github.com/gin-gonic/gin"
)

type WsService struct {
}

func (*WsService) Connect(ctx *gin.Context) {
	conn, err := ws.App(ctx.Writer, ctx.Request)
	if err != nil {
		log.C(ctx).Errorw("WebSocket connection failed")
		return
	}

	log.C(ctx).Infow("WebSocket connection established")
	code.OK.ToJson(ctx)

	id := ctx.GetInt64(known.XIdKey)

	wsClient := client.NewClient(id, conn)
	client.Manager.Register <- wsClient

	go wsClient.Read()
	go wsClient.Write()
}

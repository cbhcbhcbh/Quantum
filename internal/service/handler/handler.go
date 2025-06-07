package handler

import (
	"github.com/cbhcbhcbh/Quantum/internal/pkg/code"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
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
	defer conn.Close()

	log.C(ctx).Infow("WebSocket connection established")
	code.OK.ToJson(ctx)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.C(ctx).Errorw("Error reading message", "error", err)
			break
		}

		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.C(ctx).Errorw("Error writing message", "error", err)
			break
		}
		log.C(ctx).Infow("Message received and echoed back", "message", string(p))
	}
}

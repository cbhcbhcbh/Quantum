package router

import (
	"github.com/cbhcbhcbh/Quantum/internal/pkg/code"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	"github.com/cbhcbhcbh/Quantum/internal/service/handler"
	"github.com/gin-gonic/gin"
)

func RegisterWsRouters(engine *gin.Engine) error {
	WsService := new(handler.WsService)
	engine.NoRoute(func(c *gin.Context) {
		code.ErrPageNotFound.ToJson(c)
	})

	engine.GET("/health", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		code.OK.ToJson(c)
	})

	// websocket routers
	ws := engine.Group("/im")
	{
		ws.GET("/connect", WsService.Connect)
	}

	return nil
}

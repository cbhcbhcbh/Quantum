package router

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/controller/users"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/store"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/code"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	"github.com/gin-gonic/gin"
)

func RegisterApiRouters(engine *gin.Engine) error {
	engine.NoRoute(func(c *gin.Context) {
		code.ErrPageNotFound.ToJson(c)
	})

	engine.GET("/health", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		code.OK.ToJson(c)
	})

	uc := users.New(store.S)

	// api routers
	api := engine.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", uc.Login)
			authGroup.POST("/registered", uc.Registered)
		}
	}

	return nil
}

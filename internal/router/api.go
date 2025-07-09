package router

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/controller/friends"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/controller/group"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/controller/message"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/controller/session"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/controller/users"
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	"github.com/cbhcbhcbh/Quantum/internal/middleware"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/code"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	"github.com/gin-gonic/gin"
)

func RegisterApiRouters(engine *gin.Engine) error {

	engine.Use(middleware.Cors())

	engine.NoRoute(func(c *gin.Context) {
		code.ErrPageNotFound.ToJson(c)
	})

	engine.GET("/health", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		code.OK.ToJson(c)
	})

	uc := users.New(store.S)
	fc := friends.New(store.S)
	gc := group.New(store.S)
	sc := session.New(store.S)
	mc := message.New(store.S)

	api := engine.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", uc.Login)
			authGroup.POST("/registered", uc.Registered)
			authGroup.POST("/sendEmail", uc.SendEmail)
		}

		api.Use(middleware.Auth())
		{
			api.GET("/user/:id", uc.Info)
			api.GET("/group/list", gc.AddressList)

			api.GET("/sessions", sc.Index)
			api.POST("/sessions", sc.Store)
			api.GET("/sessions/:id", sc.Update)
			api.GET("/sessions/:id", sc.Delete)

			api.GET("/friends", fc.Index)
			api.GET("/friends/:id", fc.Show)
			api.DELETE("/friends/:id", fc.Delete)
			api.GET("/friends/status/:id", fc.GetUserStatus)

			api.POST("/friends/record", fc.SendFriendRequest)
			api.GET("/friends/record", fc.ListFriendRequests)
			api.PUT("/friends/record", fc.AcceptFriendRequest)

			// TODO: Implement message endpoints
			api.GET("/messages", mc.PrivateMessage)
			api.GET("/messages/groups", mc.GroupMessage)

			api.POST("/messages/private", mc.SendMessage)
			api.POST("/messages/group", mc.SendMessage)
			api.POST("/messages/video", mc.SendVideoMessage)
			api.POST("/messages/recall", mc.RecallMessage)

			// TODO: Implement group endpoints

			// TODO: Implement file upload endpoint
		}
	}

	return nil
}

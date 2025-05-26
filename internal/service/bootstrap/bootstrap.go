package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cbhcbhcbh/Quantum/internal/config"
	"github.com/cbhcbhcbh/Quantum/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func LoadConfig() {
	config.InitConfig()
}

func Start() error {
	gin.SetMode(viper.GetString("runmode"))
	engine := gin.Default()

	mws := []gin.HandlerFunc{gin.Recovery()}
	engine.Use(mws...)

	if err := setRoute(engine); err != nil {
		return err
	}

	httpsrv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: engine,
	}

	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return
		}
	}()

	return gracefulShutdown(httpsrv)
}

func setRoute(engine *gin.Engine) error {
	err := router.RegisterWsRouters(engine)
	if err != nil {
		return err
	}

	return nil
}

func gracefulShutdown(httpsrv *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("shutdown-timeout"))
	defer cancel()

	if err := httpsrv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

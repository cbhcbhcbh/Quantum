package quantum

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewQuantumCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quantum",
		Short: "Quantum is a CLI application",
		Long:  `A CLI application for quantum computing tasks`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	return cmd
}

func run() error {
	if err := initStore(); err != nil {
		return err
	}

	gin.SetMode(viper.GetString("runmode"))
	engine := gin.Default()

	mws := []gin.HandlerFunc{gin.Recovery()}
	engine.Use(mws...)

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

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/cbhcbhcbh/Quantum/internal/wire"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "web server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := wire.InitializeWebServer("web")
		if err != nil {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			logger.Error("initialize user server failed", zap.Error(err))
			os.Exit(1)
		}
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}

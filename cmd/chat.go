package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/cbhcbhcbh/Quantum/internal/wire"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "chat server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := wire.InitializeChatServer("chat")
		if err != nil {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			logger.Error("initialize chat server failed", zap.Error(err))
			os.Exit(1)
		}
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}

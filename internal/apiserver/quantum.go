package apiserver

import (
	"fmt"

	"github.com/cbhcbhcbh/Quantum/internal/config"
	"github.com/cbhcbhcbh/Quantum/internal/service/bootstrap"
	"github.com/spf13/cobra"
)

func NewQuantumCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quantum",
		Short: "Quantum is a CLI application",
		Long:  `A CLI application for quantum computing tasks`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			config.InitRedis()

			if err := config.InitStore(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return bootstrap.Start()
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

	cobra.OnInitialize(bootstrap.LoadConfig)

	return cmd
}

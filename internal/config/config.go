package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = "quantum.yaml"
)

func InitConfig() {
	home, err := os.Getwd()
	cobra.CheckErr(err)

	viper.AddConfigPath(filepath.Join(home, "configs"))
	viper.SetConfigType("yaml")
	viper.SetConfigName(defaultConfigName)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("QUANTUM")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: Could not read config file: %v\n", err)
		return
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
}

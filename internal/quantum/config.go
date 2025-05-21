package quantum

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cbhcbhcbh/Quantum/internal/quantum/store"
	"github.com/cbhcbhcbh/Quantum/pkg/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = "quantum.yaml"
)

func initConfig() {
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

func initStore() error {
	dbOptions := &db.PostgresOptions{
		Host:                  viper.GetString("db.host"),
		Username:              viper.GetString("db.username"),
		Password:              viper.GetString("db.password"),
		Database:              viper.GetString("db.database"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
		LogLevel:              viper.GetInt("db.log-level"),
	}

	ins, err := db.NewPostgres(dbOptions)
	if err != nil {
		return fmt.Errorf("failed to create database instance: %w", err)
	}

	_ = store.NewStore(ins)

	return nil
}

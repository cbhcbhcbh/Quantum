package config

import (
	"fmt"

	"github.com/cbhcbhcbh/Quantum/internal/apiserver/store"
	"github.com/cbhcbhcbh/Quantum/pkg/db"
	"github.com/spf13/viper"
)

func InitStore() error {
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

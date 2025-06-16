package config

import (
	"github.com/spf13/viper"
)

type EmailConfig struct {
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

func InitEmailConfig() *EmailConfig {
	mailConfig := &EmailConfig{
		Name:     viper.GetString("mail.name"),
		Password: viper.GetString("mail.password"),
		Host:     viper.GetString("mail.host"),
		Port:     viper.GetInt("mail.port"),
	}

	return mailConfig
}

package config

import (
	"github.com/spf13/viper"
)

type PostgresConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	DBName   string `mapstructure:"DB_NAME"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	SslMode  string `mapstructure:"DB_SSLMODE"`
	TimeZone string `mapstructure:"DB_TIMEZONE"`
}

func LoadPostgresConfig() PostgresConfig {
	InitEnv()

	var cfg = PostgresConfig{
		Host:     viper.GetString("DB_HOST"),
		DBName:   viper.GetString("DB_NAME"),
		Port:     viper.GetString("DB_PORT"),
		User:     viper.GetString("DB_USERNAME"),
		Password: viper.GetString("DB_PASSWORD"),
		SslMode:  viper.GetString("DB_SSLMODE"),
		TimeZone: viper.GetString("DB_TIMEZONE"),
	}

	return cfg
}

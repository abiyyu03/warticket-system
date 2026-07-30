package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresConfig
	RedisConfig
}

func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %s", err.Error())
	}

	viper.AutomaticEnv()
}

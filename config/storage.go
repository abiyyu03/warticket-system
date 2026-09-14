package config

import "github.com/spf13/viper"

type StorageConfig struct {
	Endpoint  string `mapstructure:"MINIO_ENDPOINT"`
	AccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	SecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	UseSSL    bool   `mapstructure:"MINIO_USE_SSL"`
	Bucket    string `mapstructure:"MINIO_BUCKET"`
}

func LoadStorageConfig() StorageConfig {
	InitEnv()

	var cfg = StorageConfig{
		Endpoint:  viper.GetString("MINIO_ENDPOINT"),
		AccessKey: viper.GetString("MINIO_ACCESS_KEY"),
		SecretKey: viper.GetString("MINIO_SECRET_KEY"),
		UseSSL:    viper.GetBool("MINIO_USE_SSL"),
		Bucket:    viper.GetString("MINIO_BUCKET"),
	}

	return cfg
}

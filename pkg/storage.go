package pkg

import (
	"context"
	"fmt"
	"time"

	"go-projects/hexagonal-example/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	Client *minio.Client
	Config config.StorageConfig
}

func InitStorage(cfg config.StorageConfig) (*minio.Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to minio: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check minio bucket: %w", err)
	}

	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create minio bucket: %w", err)
		}
	}

	return client, nil
}

func NewStorage(cfg config.StorageConfig) (*Storage, error) {
	client, err := InitStorage(cfg)
	if err != nil {
		return nil, err
	}

	return &Storage{
		Client: client,
		Config: cfg,
	}, nil
}

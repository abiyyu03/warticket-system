package pkg

import (
	"fmt"
	"go-projects/hexagonal-example/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SQL struct {
	*gorm.DB
	Config config.PostgresConfig
}

func InitPostgres(cfg config.PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SslMode, cfg.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.Exec("SELECT 1").Error; err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewSQL(cfg config.PostgresConfig) (*SQL, error) {
	db, err := InitPostgres(cfg)
	if err != nil {
		return nil, err
	}

	return &SQL{
		DB:     db,
		Config: cfg,
	}, nil
}

package pkg

import (
	"go-projects/hexagonal-example/config"
)

type Package struct {
	DB      *SQL
	Cache   *Redis
	Log     *Logger
	Storage *Storage
}

func NewPackage() (Package, error) {
	config.InitEnv()

	logger, err := NewLogger()
	if err != nil {
		return Package{}, err
	}

	var (
		postgresCfg = config.LoadPostgresConfig()
		redisCfg    = config.LoadRedisConfig()
		storageCfg  = config.LoadStorageConfig()
	)

	sql, err := NewSQL(postgresCfg)
	if err != nil {
		return Package{}, err
	}

	cache, err := NewRedis(redisCfg)
	if err != nil {
		return Package{}, err
	}

	storage, err := NewStorage(storageCfg)
	if err != nil {
		return Package{}, err
	}

	return Package{
		DB:      sql,
		Cache:   cache,
		Log:     logger,
		Storage: storage,
	}, nil
}

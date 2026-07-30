package pkg

import (
	"go-projects/hexagonal-example/config"
)

type Package struct {
	DB    *SQL
	Cache *Redis
}

func NewPackage() (Package, error) {
	config.InitEnv()

	var (
		postgresCfg = config.LoadPostgresConfig()
		redisCfg    = config.LoadRedisConfig()
	)

	sql, err := NewSQL(postgresCfg)
	if err != nil {
		return Package{}, err
	}

	cache, err := NewRedis(redisCfg)
	if err != nil {
		return Package{}, err
	}

	return Package{
		DB:    sql,
		Cache: cache,
	}, nil
}

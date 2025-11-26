package cache

import (
	"go-projects/hexagonal-example/internal/adapter/outbound/cache/user"

	"go.uber.org/dig"
)

type Cache struct {
	dig.In

	User user.Cache
}

func Register(container *dig.Container) error {
	if err := container.Provide(user.New); err != nil {
		return err
	}

	return nil
}

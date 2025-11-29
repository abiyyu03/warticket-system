package cache

import (
	"go-projects/hexagonal-example/internal/adapter/outbound/cache/ticket"
	"go-projects/hexagonal-example/internal/adapter/outbound/cache/user"

	"go.uber.org/dig"
)

type Cache struct {
	dig.In

	User   user.Cache
	Ticket ticket.Cache
}

func Register(container *dig.Container) error {
	if err := container.Provide(user.New); err != nil {
		return err
	}
	if err := container.Provide(ticket.New); err != nil {
		return err
	}

	return nil
}

package repository

import (
	"go-projects/hexagonal-example/internal/adapter/outbound/repository/event"
	"go-projects/hexagonal-example/internal/adapter/outbound/repository/user"
	"go-projects/hexagonal-example/internal/adapter/outbound/repository/userTicket"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In

	User       user.Repository
	UserTicket userTicket.Repository
	Event      event.Repository
}

func Register(container *dig.Container) error {
	if err := container.Provide(user.New); err != nil {
		return err
	}
	if err := container.Provide(userTicket.New); err != nil {
		return err
	}
	if err := container.Provide(event.New); err != nil {
		return err
	}

	return nil
}

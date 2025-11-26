package service

import (
	"go-projects/hexagonal-example/internal/service/ticket"
	"go-projects/hexagonal-example/internal/service/user"

	"go.uber.org/dig"
)

type Service struct {
	dig.In

	User   user.UserService
	Ticket ticket.TicketService
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

package service

import (
	"go-projects/hexagonal-example/internal/service/author"
	"go-projects/hexagonal-example/internal/service/event"
	"go-projects/hexagonal-example/internal/service/ticket"
	"go-projects/hexagonal-example/internal/service/user"

	"go.uber.org/dig"
)

type Service struct {
	dig.In

	User   user.UserService
	Event  event.EventService
	Ticket ticket.TicketService
	Author author.AuthorService
}

func Register(container *dig.Container) error {
	if err := container.Provide(user.New); err != nil {
		return err
	}

	if err := container.Provide(ticket.New); err != nil {
		return err
	}

	if err := container.Provide(event.New); err != nil {
		return err
	}

	if err := container.Provide(author.New); err != nil {
		return err
	}

	return nil
}

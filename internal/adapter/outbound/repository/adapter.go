package repository

import (
	"go-projects/hexagonal-example/internal/adapter/outbound/repository/event"
	"go-projects/hexagonal-example/internal/adapter/outbound/repository/eventDetail"
	"go-projects/hexagonal-example/internal/adapter/outbound/repository/seat"
	"go-projects/hexagonal-example/internal/adapter/outbound/repository/user"
	"go-projects/hexagonal-example/internal/adapter/outbound/repository/userEventTicket"

	"go.uber.org/dig"
)

type Repository struct {
	dig.In

	User            user.Repository
	UserEventTicket userEventTicket.Repository
	Event           event.Repository
	Seat            seat.Repository
	EventDetail     eventDetail.Repository
}

func Register(container *dig.Container) error {
	if err := container.Provide(user.New); err != nil {
		return err
	}
	if err := container.Provide(userEventTicket.New); err != nil {
		return err
	}
	if err := container.Provide(event.New); err != nil {
		return err
	}
	if err := container.Provide(seat.New); err != nil {
		return err
	}
	if err := container.Provide(eventDetail.New); err != nil {
		return err
	}

	return nil
}

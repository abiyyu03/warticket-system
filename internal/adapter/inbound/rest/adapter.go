package rest

import (
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/ticket"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/user"

	"go-projects/hexagonal-example/internal/adapter/inbound/rest/event"

	"go.uber.org/dig"
)

type Inbound struct {
	dig.In

	User   user.Handler
	Ticket ticket.Handler
	Event  event.Handler
}

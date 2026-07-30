package user

import (
	"go-projects/hexagonal-example/internal/service"

	"go.uber.org/dig"
)

type Handler struct {
	dig.In

	Service service.Service
}

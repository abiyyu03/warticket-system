package di

import (
	"go-projects/hexagonal-example/internal/adapter/outbound"
	"go-projects/hexagonal-example/internal/service"
	"go-projects/hexagonal-example/pkg"

	"go.uber.org/dig"
)

func Container() (*dig.Container, error) {
	container := dig.New()

	if err := outbound.Register(container); err != nil {
		return nil, err
	}

	if err := pkg.Register(container); err != nil {
		return nil, err
	}

	if err := service.Register(container); err != nil {
		return nil, err
	}

	return container, nil
}

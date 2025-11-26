package repository

import (
	"go.uber.org/dig"
)

type External struct {
	dig.In

	// inject external service here
	// Mail mail
	// Payment payment
}

func Register(container *dig.Container) error {
	// if err := container.Provide(user.New); err != nil {
	// 	return err
	// }

	return nil
}

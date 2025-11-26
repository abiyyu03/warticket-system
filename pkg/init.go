package pkg

import (
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	if err := container.Provide(NewPackage); err != nil {
		return err
	}

	return nil
}

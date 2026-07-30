package main

import (
	"go-projects/hexagonal-example/internal/adapter/inbound/rest"
	"go-projects/hexagonal-example/pkg"
	"go-projects/hexagonal-example/pkg/di"

	"github.com/gofiber/fiber/v2"
)

func main() {
	var app = fiber.New()

	container, err := di.Container()
	if err != nil {
		panic(err)
	}

	err = container.Invoke(func(pkg pkg.Package, inbound rest.Inbound) error {
		inbound.ApiRoutes(app)

		app.Listen(":8080")

		return nil
	})
	if err != nil {
		panic(err)
	}

}

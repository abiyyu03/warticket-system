package rest

import "github.com/gofiber/fiber/v2"

func (i Inbound) ApiRoutes(app *fiber.App) {
	app.Get("/health-check", func(ctx *fiber.Ctx) error {
		return ctx.SendString("OK LURD")
	})

	v1Api := app.Group("v1/api")
	user := v1Api.Group("users")
	user.Get("/", i.User.GetAll)

	v1Api.Post("/register", i.User.RegisterUser)

	ticket := v1Api.Group("tickets")
	ticket.Post("/init-order", i.Ticket.InitOrder)
	ticket.Post("/claim", i.Ticket.Purchase)
	ticket.Post("/redeem", i.Ticket.Redeem)

	event := v1Api.Group("events")
	event.Get("/", i.Event.GetListEvent)
	event.Post("/", i.Event.CreateEvent)
	event.Get("/:id/form", i.Event.GetEventForm)
}

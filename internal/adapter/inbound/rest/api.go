package rest

import (
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/middleware"
	"go-projects/hexagonal-example/pkg/token"

	"github.com/gofiber/fiber/v2"
)

func (i Inbound) ApiRoutes(app *fiber.App) {
	app.Get("/health-check", func(ctx *fiber.Ctx) error {
		return ctx.SendString("OK LURD")
	})

	v1Api := app.Group("v1/api")
	user := v1Api.Group("users")
	user.Get("/", i.User.GetAll)
	v1Api.Post("/register", i.User.RegisterUser)

	// user guest event
	event := v1Api.Group("events")
	event.Get("/", i.Event.GetListEvent)
	event.Get("/:id", i.Event.GetOneEvent)

	// user ticket
	ticket := v1Api.Group("tickets")
	ticket.Post("/init-order", i.Ticket.InitOrder)
	ticket.Post("/claim", i.Ticket.Purchase)
	ticket.Post("/redeem", i.Ticket.Redeem)

	// author area
	author := v1Api.Group("authors")

	// --- auth publik (tanpa token) ---
	author.Post("/register", i.Author.Register)
	author.Post("/login", i.Author.Login)
	author.Post("/refresh", i.Author.Refresh)

	// --- mulai di sini wajib access token role author/admin ---
	author.Use(middleware.Auth(token.RoleAuthor, token.RoleAdmin))

	// pengelolaan user: pakai ulang daftar user yang sudah ada.
	author.Get("/users", i.User.GetAll)

	authorEvent := author.Group("events")
	authorEvent.Get("/", i.Event.GetListEvent)
	authorEvent.Get("/:id", i.Event.GetOneEvent)
	authorEvent.Post("/", i.Event.CreateEvent)
	authorEvent.Get("/:id/form", i.Event.GetEventForm)
	authorEvent.Get("/:id/tickets", i.Ticket.GetEventTickets)
	authorEvent.Get("/:id/tickets/export", i.Ticket.ExportEventTickets)
}

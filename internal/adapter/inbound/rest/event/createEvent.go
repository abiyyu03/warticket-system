package ticket

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/event/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateEvent(fctx *fiber.Ctx) error {
	var (
		request entity.CreateEventRequest
		ctx     = context.WithValue(fctx.Context(), "x-user-id", fctx.Get("x-user-id"))
	)

	err := h.Service.Event.CreateEvent(ctx, request.ToUcEntity())
	if err != nil {
		return err
	}

	return fctx.Status(fiber.StatusOK).JSON(
		baseEntity.BaseResponse{}.ToResponse(
			"Ticket Purchase Successfully",
			fiber.StatusCreated,
			nil,
			nil,
		),
	)
}

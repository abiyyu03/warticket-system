package ticket

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/ticket/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) InitOrder(fctx *fiber.Ctx) error {
	var (
		ctx     = context.WithValue(fctx.Context(), "x-user-id", fctx.Get("x-user-id"))
		request entity.InitOrderRequest
	)

	if err := fctx.BodyParser(&request); err != nil {
		return err
	}

	response, err := h.Service.Ticket.InitOrder(ctx, request.ToUcEntity())
	if err != nil {
		return err
	}

	return fctx.Status(fiber.StatusOK).JSON(
		baseEntity.BaseResponse{}.ToResponse(
			"Init Order Successfully",
			fiber.StatusOK,
			response,
			nil,
		),
	)
}

package ticket

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/ticket/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Purchase(fctx *fiber.Ctx) error {
	var (
		request entity.PurchaseRequest
		ctx     = context.WithValue(fctx.Context(), "x-user-id", fctx.Get("x-user-id"))
	)

	response, err := h.Service.Ticket.Purchase(ctx, request.ToUcEntity())
	if err != nil {
		return err
	}

	return fctx.Status(fiber.StatusOK).JSON(
		baseEntity.BaseResponse{}.ToResponse(
			"Ticket Purchase Successfully",
			fiber.StatusOK,
			response,
			nil,
		),
	)
}

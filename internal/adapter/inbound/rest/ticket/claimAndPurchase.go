package ticket

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/ticket/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ClaimAndPurchase(fctx *fiber.Ctx) error {
	var ctx = context.Background()
	var (
		request entity.ClaimTicketRequest
	)

	response, err := h.Service.Ticket.ClaimAndPurchase(ctx, request.ToUcEntity())
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

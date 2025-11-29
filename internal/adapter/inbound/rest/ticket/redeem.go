package ticket

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/ticket/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Redeem(fctx *fiber.Ctx) error {
	var ctx = context.Background()
	var (
		request entity.RedeemRequest
	)

	if err := fctx.BodyParser(&request); err != nil {
		return err
	}

	err := h.Service.Ticket.Redeem(ctx, request.ToUcEntity())
	if err != nil {
		return err
	}

	return fctx.Status(fiber.StatusOK).JSON(
		baseEntity.BaseResponse{}.ToResponse(
			"Redeem Successfully",
			fiber.StatusOK,
			nil,
			nil,
		),
	)
}

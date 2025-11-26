package ticket

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/ticket/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Claim(fctx *fiber.Ctx) error {
	var ctx = context.Background()
	var (
		request entity.ClaimTicketRequest
	)

	err := h.Service.Ticket.Claim(ctx, request.ToUcEntity())
	if err != nil {
		return err
	}

	return nil

	// return fctx.Status(fiber.StatusOK).JSON(
	// 	response.ToResponse(
	// 		"User fetched successfully",
	// 		fiber.StatusOK,
	// 		users,
	// 		nil,
	// 	),
	// )
}

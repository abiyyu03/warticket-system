package ticket

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetEventTickets(fctx *fiber.Ctx) error {
	ctx := context.Background()

	eventID, err := strconv.ParseInt(fctx.Params("id"), 10, 64)
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).JSON(
			baseEntity.BaseResponse{}.ToResponse("event id tidak valid", fiber.StatusBadRequest, nil, nil),
		)
	}

	tickets, err := h.Service.Ticket.ListEventTickets(ctx, eventID)
	if err != nil {
		return err
	}

	return fctx.Status(fiber.StatusOK).JSON(
		baseEntity.BaseResponse{}.ToResponse("Event Tickets Retrieved Successfully", fiber.StatusOK, tickets, nil),
	)
}

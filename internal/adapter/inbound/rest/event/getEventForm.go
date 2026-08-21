package event

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetEventForm(fctx *fiber.Ctx) error {
	ctx := context.WithValue(fctx.Context(), "x-user-id", fctx.Get("x-user-id"))

	eventID, err := strconv.ParseInt(fctx.Params("id"), 10, 64)
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).JSON(
			baseEntity.BaseResponse{}.ToResponse("event id tidak valid", fiber.StatusBadRequest, nil, nil),
		)
	}

	form, err := h.Service.Event.GetEventForm(ctx, eventID)
	if err != nil {
		return err
	}

	return fctx.Status(fiber.StatusOK).JSON(
		baseEntity.BaseResponse{}.ToResponse("Event Form Retrieved Successfully", fiber.StatusOK, form, nil),
	)
}

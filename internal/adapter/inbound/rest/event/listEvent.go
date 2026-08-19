package event

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/event/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetListEvent(fctx *fiber.Ctx) error {
	var (
		request entity.GetListEventRequest
		ctx     = context.WithValue(fctx.Context(), "x-user-id", fctx.Get("x-user-id"))
	)

	resp, err := h.Service.Event.GetListEvent(ctx, request.ToUcEntity())
	if err != nil {
		return err
	}

	return fctx.Status(fiber.StatusOK).JSON(
		baseEntity.BaseResponse{}.ToResponse(
			"Event Fetched Successfully",
			fiber.StatusOK,
			resp,
			nil,
		),
	)
}

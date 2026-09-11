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
		ctx     = context.WithValue(fctx.UserContext(), "x-user-id", fctx.Get("x-user-id"))
	)

	if err := fctx.BodyParser(&request); err != nil {
		return err
	}

	response, err := h.Service.Ticket.Purchase(ctx, request.ToUcEntity())
	if err != nil {
		// mayoritas error di checkout bersifat validasi/bisnis (email/opsi form,
		// tx_id tidak valid) -> balas 400 dengan pesan jelas.
		return fctx.Status(fiber.StatusBadRequest).JSON(
			baseEntity.BaseResponse{}.ToResponse(err.Error(), fiber.StatusBadRequest, nil, nil),
		)
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

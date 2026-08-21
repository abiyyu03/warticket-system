package event

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/event/entity"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SubmitRegistration(fctx *fiber.Ctx) error {
	ctx := context.WithValue(fctx.Context(), "x-user-id", fctx.Get("x-user-id"))

	eventID, err := strconv.ParseInt(fctx.Params("id"), 10, 64)
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).JSON(
			baseEntity.BaseResponse{}.ToResponse("event id tidak valid", fiber.StatusBadRequest, nil, nil),
		)
	}

	var request entity.SubmitRegistrationRequest
	if err := fctx.BodyParser(&request); err != nil {
		return err
	}

	if err := h.Service.Event.SubmitRegistration(ctx, request.ToUcEntity(eventID)); err != nil {
		// mayoritas error di sini bersifat validasi/bisnis (opsi salah, wajib
		// diisi, sudah terdaftar) -> kembalikan 400 dengan pesan jelas.
		return fctx.Status(fiber.StatusBadRequest).JSON(
			baseEntity.BaseResponse{}.ToResponse(err.Error(), fiber.StatusBadRequest, nil, nil),
		)
	}

	return fctx.Status(fiber.StatusOK).JSON(
		baseEntity.BaseResponse{}.ToResponse("Registration Submitted Successfully", fiber.StatusCreated, nil, nil),
	)
}

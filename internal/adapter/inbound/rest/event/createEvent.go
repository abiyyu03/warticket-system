package event

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/event/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateEvent(fctx *fiber.Ctx) error {
	var (
		request entity.CreateEventRequest
		ctx     = context.WithValue(fctx.Context(), "x-user-id", fctx.Get("x-user-id"))
	)

	if err := fctx.BodyParser(&request); err != nil {
		return err
	}

	// definisi formulir kustom (opsional) dikirim sebagai JSON di multipart.
	formFields, err := request.ParseFormFields()
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).JSON(
			baseEntity.BaseResponse{}.ToResponse("form_fields tidak valid", fiber.StatusBadRequest, nil, nil),
		)
	}

	ucRequest := request.ToUcEntity()
	ucRequest.FormFields = formFields

	if err := h.Service.Event.CreateEvent(ctx, ucRequest); err != nil {
		return err
	}

	return fctx.Status(fiber.StatusOK).JSON(
		baseEntity.BaseResponse{}.ToResponse(
			"Event Created Successfully",
			fiber.StatusCreated,
			nil,
			nil,
		),
	)
}

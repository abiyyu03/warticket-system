package author

import (
	"context"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/author/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Register(fctx *fiber.Ctx) error {
	ctx := context.Background()

	var request entity.RegisterRequest
	if err := fctx.BodyParser(&request); err != nil {
		return err
	}

	if err := h.Service.Author.Register(ctx, request.ToUcEntity()); err != nil {
		return fctx.Status(fiber.StatusBadRequest).JSON(
			baseEntity.BaseResponse{}.ToResponse(err.Error(), fiber.StatusBadRequest, nil, nil),
		)
	}

	return fctx.Status(fiber.StatusCreated).JSON(
		baseEntity.BaseResponse{}.ToResponse("Author Registered Successfully", fiber.StatusCreated, nil, nil),
	)
}

package author

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/author/entity"
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Login(fctx *fiber.Ctx) error {
	ctx := context.Background()

	var request entity.LoginRequest
	if err := fctx.BodyParser(&request); err != nil {
		return err
	}

	resp, err := h.Service.Author.Login(ctx, request.ToUcEntity())
	if err != nil {
		return fctx.Status(fiber.StatusUnauthorized).JSON(
			baseEntity.BaseResponse{}.ToResponse(err.Error(), fiber.StatusUnauthorized, nil, nil),
		)
	}

	return fctx.Status(fiber.StatusOK).JSON(
		baseEntity.BaseResponse{}.ToResponse("Login Successfully", fiber.StatusOK, resp, nil),
	)
}

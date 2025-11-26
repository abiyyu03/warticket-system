package user

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/inbound/rest/user/entity"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetAll(fctx *fiber.Ctx) error {
	var ctx = context.Background()
	var (
		response entity.BaseResponse
	)

	users, err := h.Service.User.GetAll(ctx)
	if err != nil {
		return err
	}

	return fctx.Status(fiber.StatusOK).JSON(
		response.ToResponse(
			"User fetched successfully",
			fiber.StatusOK,
			users,
			nil,
		),
	)
}

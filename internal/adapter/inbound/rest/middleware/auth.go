package middleware

import (
	baseEntity "go-projects/hexagonal-example/internal/adapter/inbound/rest/entity"
	"go-projects/hexagonal-example/pkg/token"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const bearerPrefix = "Bearer "

// Auth memvalidasi access token JWT dan (opsional) membatasi role yang diizinkan.
// Tanpa argumen role, cukup token valid. Jika role diberikan, role pada token
// harus termasuk salah satunya.
//
// Efek samping: menaruh user id & role di Locals dan menimpa header x-user-id
// dengan id dari token, supaya handler lama yang membaca x-user-id tetap jalan
// memakai identitas terautentikasi.
func Auth(allowedRoles ...string) fiber.Handler {
	manager := token.New()

	return func(c *fiber.Ctx) error {
		authz := c.Get(fiber.HeaderAuthorization)
		if !strings.HasPrefix(authz, bearerPrefix) {
			return unauthorized(c, "authorization header tidak ada atau salah format")
		}

		claims, err := manager.ParseAccess(strings.TrimPrefix(authz, bearerPrefix))
		if err != nil {
			return unauthorized(c, "token tidak valid atau kedaluwarsa")
		}

		if len(allowedRoles) > 0 && !contains(allowedRoles, claims.Role) {
			return forbidden(c, "role tidak diizinkan mengakses resource ini")
		}

		c.Locals("user-id", claims.UserID)
		c.Locals("role", claims.Role)
		c.Request().Header.Set("x-user-id", strconv.FormatInt(claims.UserID, 10))
		return c.Next()
	}
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func unauthorized(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(
		baseEntity.BaseResponse{}.ToResponse(msg, fiber.StatusUnauthorized, nil, nil),
	)
}

func forbidden(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusForbidden).JSON(
		baseEntity.BaseResponse{}.ToResponse(msg, fiber.StatusForbidden, nil, nil),
	)
}

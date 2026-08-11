package ticket

import (
	"context"
	entity "go-projects/hexagonal-example/internal/service/entity/ticket"
)

func (s service) Redeem(ctx context.Context, redeem entity.RedeemRequest) error {
	// var (
	// 	orm      = s.repository.DB
	// 	redeemAt = time.Now()
	// )
	// redeem by update status

	return nil
}

type IRedeem interface {
	Redeem(ctx context.Context, user entity.RedeemRequest) error
}

package ticket

import (
	"context"
	entity "go-projects/hexagonal-example/internal/service/entity/ticket"
	"time"
)

func (s service) Redeem(ctx context.Context, redeem entity.RedeemRequest) error {
	var (
		orm      = s.repository.DB
		redeemAt = time.Now()
	)
	// redeem by update status
	err := s.Repository.UserEventTicket.UpdateByCode(ctx, orm, redeem.ToObEntityRedeem(&redeemAt))
	if err != nil {
		return err
	}

	return nil
}

type IRedeem interface {
	Redeem(ctx context.Context, user entity.RedeemRequest) error
}

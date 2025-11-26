package ticket

import (
	"context"
	entity "go-projects/hexagonal-example/internal/service/entity/user"
)

func (s service) Redeem(ctx context.Context, user entity.User) error {
	// err := s.Repository.User.Redeem(ctx, user.ToOutbondUser())
	// if err != nil {
	// 	return err
	// }

	return nil
}

type IRedeem interface {
	Redeem(ctx context.Context, user entity.User) error
}

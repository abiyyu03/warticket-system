package user

import (
	"context"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/user"
)

func (s service) CreateUser(ctx context.Context, user ucEntity.User) error {
	err := s.Repository.User.CreateUser(ctx, user.ToOutbondUser())
	if err != nil {
		return err
	}

	return nil
}

type ICreateUser interface {
	CreateUser(ctx context.Context, user ucEntity.User) error
}

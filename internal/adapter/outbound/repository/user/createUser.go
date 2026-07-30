package user

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

func (r user) CreateUser(ctx context.Context, user entity.User) error {
	if err := r.Package.DB.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	return nil
}

type ICreateUser interface {
	CreateUser(ctx context.Context, user entity.User) error
}

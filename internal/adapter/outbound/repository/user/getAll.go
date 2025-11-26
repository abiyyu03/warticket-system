package user

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

func (r user) GetAll(ctx context.Context) ([]entity.User, error) {
	var user []entity.User
	if err := r.Package.DB.WithContext(ctx).Find(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

type IGetAll interface {
	GetAll(ctx context.Context) ([]entity.User, error)
}

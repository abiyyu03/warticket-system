package user

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

func (r user) GetUser(ctx context.Context) (entity.User, error) {
	var response entity.User
	if err := r.Package.DB.WithContext(ctx).Find(&response).Error; err != nil {
		return response, err
	}

	return response, nil
}

type IGetUser interface {
	GetUser(ctx context.Context) (entity.User, error)
}

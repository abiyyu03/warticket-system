package user

import (
	"context"
	"encoding/json"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

func (c userCache) GetAllUser(ctx context.Context) ([]entity.User, error) {
	var key = "users:getAll"
	result := c.Package.Cache.Client.Get(ctx, key)

	var users []entity.User
	str, err := result.Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(str), &users); err != nil {
		return nil, err
	}

	return users, nil
}

type IGetAllUser interface {
	GetAllUser(ctx context.Context) ([]entity.User, error)
}

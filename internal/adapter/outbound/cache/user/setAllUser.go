package user

import (
	"context"
	"encoding/json"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"time"
)

func (c userCache) SetAllUser(ctx context.Context, user []entity.User) error {
	var key = "users:getAll"

	users, err := json.Marshal(&user)
	if err != nil {
		return err
	}

	result := c.Package.Cache.Client.Set(ctx, key, users, time.Minute*10)

	_, err = result.Result()
	if err != nil {
		return err
	}

	return nil
}

type ISetAllUser interface {
	SetAllUser(ctx context.Context, user []entity.User) error
}

package ticket

import (
	"context"
	"encoding/json"
	"fmt"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

func (c ticketCache) GetInitOrder(ctx context.Context, req entity.CacheInitOrderRequest) (entity.CacheInitOrderResponse, error) {
	redisKey := fmt.Sprintf(key, req.UserID)
	result := c.Package.Cache.Client.Get(ctx, redisKey)

	var initUserOrder entity.CacheInitOrderResponse
	str, err := result.Result()
	if err != nil {
		return initUserOrder, err
	}

	if err := json.Unmarshal([]byte(str), &initUserOrder); err != nil {
		return initUserOrder, err
	}

	return initUserOrder, nil
}

type IGetInitOrder interface {
	GetInitOrder(ctx context.Context, req entity.CacheInitOrderRequest) (entity.CacheInitOrderResponse, error)
}

package ticket

import (
	"context"
	"encoding/json"
	"fmt"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"time"
)

func (c ticketCache) SetInitOrder(ctx context.Context, req entity.CacheInitOrderRequest) error {
	init, err := json.Marshal(&entity.CacheInitOrderRequest{
		Date:     req.Date,
		EventID:  req.EventID,
		Quantity: req.Quantity,
		UserID:   req.UserID,
		Price:    req.Price,
	})
	if err != nil {
		return err
	}

	redisKey := fmt.Sprintf(key, req.UserID, req.EventID)

	result := c.Package.Cache.Client.Set(ctx, redisKey, init, time.Minute*10)

	_, err = result.Result()
	if err != nil {
		return err
	}

	return nil
}

type ISetInitOrder interface {
	SetInitOrder(ctx context.Context, req entity.CacheInitOrderRequest) error
}

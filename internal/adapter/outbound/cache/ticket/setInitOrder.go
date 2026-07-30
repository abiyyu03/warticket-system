package ticket

import (
	"context"
	"encoding/json"
	"fmt"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"time"
)

func (c ticketCache) SetInitOrder(ctx context.Context, req entity.CacheInitOrderResponse) error {
	init, err := json.Marshal(&entity.CacheInitOrderResponse{
		ChairCode: req.ChairCode,
		Date:      req.Date,
		EventID:   req.EventID,
		UserID:    req.UserID,
	})
	if err != nil {
		return err
	}

	redisKey := fmt.Sprintf(key, req.UserID)

	result := c.Package.Cache.Client.Set(ctx, redisKey, init, time.Minute*10)

	_, err = result.Result()
	if err != nil {
		return err
	}

	return nil
}

type ISetInitOrder interface {
	SetInitOrder(ctx context.Context, req entity.CacheInitOrderResponse) error
}

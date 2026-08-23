package ticket

import (
	"context"
	"fmt"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

func (c ticketCache) ClearInitOrder(ctx context.Context, req entity.CacheInitOrderRequest) error {
	redisKey := fmt.Sprintf(key, req.TxID)
	c.Package.Cache.Client.Del(ctx, redisKey)

	return nil
}

type IClearInitOrder interface {
	ClearInitOrder(ctx context.Context, req entity.CacheInitOrderRequest) error
}

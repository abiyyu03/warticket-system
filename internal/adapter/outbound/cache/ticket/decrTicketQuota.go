package ticket

import (
	"context"
	"fmt"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

func (c ticketCache) DecrTicketQuota(ctx context.Context, req entity.DecrTicketQuotaRequest) error {
	redisKey := fmt.Sprintf(keyDecrQuota, req.EventID)

	err := c.Package.Cache.Client.Decr(ctx, redisKey)
	if err != nil {
		return nil
	}

	return nil
}

type IDecrTicketQuota interface {
	DecrTicketQuota(ctx context.Context, req entity.DecrTicketQuotaRequest) error
}

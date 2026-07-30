package ticket

import (
	"context"
	"fmt"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)
 
func (c ticketCache) SetTicketQuota(ctx context.Context, req entity.SetTicketQuotaRequest) error {
	redisKey := fmt.Sprintf(keyDecrQuota, req.EventID)

	if err := c.Package.Cache.Client.Set(ctx, redisKey, req.Quota, 0).Err(); err != nil {
		return err
	}

	return nil
}

type ISetTicketQuota interface {
	SetTicketQuota(ctx context.Context, req entity.SetTicketQuotaRequest) error
}
package ticket

import (
	"context"
	"errors"
	"fmt"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

// ErrQuotaSoldOut dikembalikan saat sisa kuota Redis tidak cukup direservasi.
var ErrQuotaSoldOut = errors.New("event sold out")

func (c ticketCache) DecrTicketQuota(ctx context.Context, req entity.DecrTicketQuotaRequest) error {
	redisKey := fmt.Sprintf(keyDecrQuota, req.EventID)

	// DecrBy atomik dan mengembalikan sisa setelah pengurangan.
	remaining, err := c.Package.Cache.Client.DecrBy(ctx, redisKey, req.Quantity).Result()
	if err != nil {
		return err
	}

	// stok tidak cukup: kembalikan (compensate) lalu tolak sebagai sold out.
	if remaining < 0 {
		c.Package.Cache.Client.IncrBy(ctx, redisKey, req.Quantity)
		return ErrQuotaSoldOut
	}

	return nil
}

type IDecrTicketQuota interface {
	DecrTicketQuota(ctx context.Context, req entity.DecrTicketQuotaRequest) error
}

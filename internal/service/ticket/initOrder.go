package ticket

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
	"strconv"
	"time"
)

func (s service) InitOrder(ctx context.Context, req ucEntity.InitOrderRequest) (ucEntity.InitOrderResponse, error) {
	var (
		orm       = s.repository.DB
		response  ucEntity.InitOrderResponse
		userIdCtx = ctx.Value("x-user-id").(string)
		userId, _ = strconv.ParseInt(userIdCtx, 10, 64)
	)

	// check cache availability if already exist
	initOrder, err := s.Cache.Ticket.GetInitOrder(ctx, entity.CacheInitOrderRequest{UserID: userId})
	if err == nil {
		return ucEntity.InitOrderResponse{
			Date:     initOrder.Date,
			EventID:  initOrder.EventID,
			Quantity: initOrder.Quantity,
		}, nil
	}

	// use REDIS DECR by quantity
	for i := 0; i < int(req.Quantity); i++ {
		err = s.Cache.Ticket.DecrTicketQuota(ctx, entity.DecrTicketQuotaRequest{EventID: initOrder.EventID})
		if err != nil {
			return response, err
		}
	}

	parsedEventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return response, err
	}

	// check event date
	_, err = s.Repository.Event.GetOneById(ctx, orm, req.ToObEvent(parsedEventDate))
	if err != nil {
		return response, err
	}

	// caching initialization data
	err = s.Cache.Ticket.SetInitOrder(ctx, req.ToObSetCache(userId))
	if err != nil {
		return response, err
	}

	return ucEntity.InitOrderResponse{
		Date:     req.Date,
		EventID:  req.EventID,
		Quantity: req.Quantity,
	}, nil
}

type IInitOrder interface {
	InitOrder(ctx context.Context, req ucEntity.InitOrderRequest) (ucEntity.InitOrderResponse, error)
}

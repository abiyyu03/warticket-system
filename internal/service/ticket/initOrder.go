package ticket

import (
	"context"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
	"slices"
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
	cachedInitOrder, err := s.Cache.Ticket.GetInitOrder(ctx, req.ToObGetCache(userId))
	if err == nil {
		// another user choose same chair
		for _, v := range cachedInitOrder.ChairCode {
			if slices.Contains(req.ChairCode, v) && cachedInitOrder.UserID == userId {
				return ucEntity.InitOrderResponse{
					Status: "BOOKED_TEMPORARY",
				}, err
			}
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

	// caching payload
	err = s.Cache.Ticket.SetInitOrder(ctx, req.ToObSetCache(userId))
	if err != nil {
		return response, err
	}

	return ucEntity.InitOrderResponse{
		Status: "INITIATED_SUCCESS",
	}, nil
}

type IInitOrder interface {
	InitOrder(ctx context.Context, req ucEntity.InitOrderRequest) (ucEntity.InitOrderResponse, error)
}

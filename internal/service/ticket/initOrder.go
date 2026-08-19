package ticket

import (
	"context"
	"errors"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	obEntity "go-projects/hexagonal-example/internal/adapter/outbound/entity"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s service) InitOrder(ctx context.Context, req ucEntity.InitOrderRequest) (ucEntity.InitOrderResponse, error) {
	var (
		orm       = s.repository.DB
		response  ucEntity.InitOrderResponse
		userIdCtx = ctx.Value("x-user-id").(string)
		userId, _ = strconv.ParseInt(userIdCtx, 10, 64)
	)

	// logger tercakup untuk seluruh flow ini; field dasar ikut di setiap baris log.
	log := s.logger.With(
		zap.String("flow", "init_order"),
		zap.Int64("user_id", userId),
		zap.Int64("event_id", req.EventID),
		zap.Int64("quantity", req.Quantity),
	)
	log.Info("init order requested")

	// check cache availability if already exist
	initOrder, err := s.Cache.Ticket.GetInitOrder(ctx, entity.CacheInitOrderRequest{UserID: userId})
	if err == nil {
		log.Info("reservation already exists, returning cached order",
			zap.String("cached_date", initOrder.Date))
		return ucEntity.InitOrderResponse{
			Date:     initOrder.Date,
			EventID:  initOrder.EventID,
			Quantity: initOrder.Quantity,
		}, nil
	}

	// validate event
	_, err = s.Repository.Event.GetOneById(ctx, orm, obEntity.Event{ID: req.EventID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("event not found")
			return response, errors.New("event not found")
		}
		log.Error("failed to load event", zap.Error(err))
		return response, err
	}

	// use REDIS DECR by quantity
	for i := 0; i < int(req.Quantity); i++ {
		err = s.Cache.Ticket.DecrTicketQuota(ctx, entity.DecrTicketQuotaRequest{EventID: req.EventID})
		if err != nil {
			log.Error("failed to decrement ticket quota in redis",
				zap.Int("at_slot", i), zap.Error(err))
			return response, err
		}
	}

	parsedEventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		log.Warn("invalid event date format", zap.String("date", req.Date), zap.Error(err))
		return response, err
	}

	// check event date
	_, err = s.Repository.Event.GetOneById(ctx, orm, req.ToObEvent(parsedEventDate))
	if err != nil {
		log.Error("event date validation failed", zap.Error(err))
		return response, err
	}

	// caching initialization data
	err = s.Cache.Ticket.SetInitOrder(ctx, req.ToObSetCache(userId))
	if err != nil {
		log.Error("failed to cache init order", zap.Error(err))
		return response, err
	}

	log.Info("init order initiated successfully")
	return ucEntity.InitOrderResponse{
		Date:     req.Date,
		EventID:  req.EventID,
		Quantity: req.Quantity,
		Price:    req.Price,
	}, nil
}

type IInitOrder interface {
	InitOrder(ctx context.Context, req ucEntity.InitOrderRequest) (ucEntity.InitOrderResponse, error)
}

package ticket

import (
	"context"
	"errors"
	cacheTicket "go-projects/hexagonal-example/internal/adapter/outbound/cache/ticket"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
	"strconv"

	"github.com/google/uuid"
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

	// tx_id diterbitkan di awal init order; dipakai sebagai kunci reservasi
	// (cache) sekaligus id transaksi saat purchase nanti.
	txID := uuid.NewString()

	log := s.logger.With(
		zap.String("flow", "init_order"),
		zap.String("tx_id", txID),
		zap.Int64("user_id", userId),
		zap.Int64("event_id", req.EventID),
		zap.Int64("quantity", req.Quantity),
	)
	log.Info("init order requested")

	// ambil event (tanpa validasi tanggal; date diisi server-side).
	event, err := s.Repository.Event.GetByID(ctx, orm, req.EventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("event not found")
			return response, errors.New("event not found")
		}
		log.Error("failed to fetch event", zap.Error(err))
		return response, err
	}

	// reserve kuota di redis: satu DECRBY sebanyak quantity (atomik).
	if err = s.Cache.Ticket.DecrTicketQuota(ctx, entity.DecrTicketQuotaRequest{
		EventID:  req.EventID,
		Quantity: req.Quantity,
	}); err != nil {
		if errors.Is(err, cacheTicket.ErrQuotaSoldOut) {
			log.Warn("event sold out, reservation rejected")
			return response, err
		}
		log.Error("failed to decrement ticket quota in redis", zap.Error(err))
		return response, err
	}

	// date diisi dari tanggal event (server-side), bukan dari client.
	date := event.StartDate.Format("2006-01-02")

	// cache reservasi di-key oleh tx_id.
	if err = s.Cache.Ticket.SetInitOrder(ctx, req.ToObSetCache(txID, date, userId, int64(event.Price))); err != nil {
		log.Error("failed to cache init order", zap.Error(err))
		return response, err
	}

	log.Info("init order initiated successfully")
	return ucEntity.InitOrderResponse{
		TxID:     txID,
		Date:     date,
		EventID:  req.EventID,
		Quantity: req.Quantity,
		Price:    int64(event.Price),
	}, nil
}

type IInitOrder interface {
	InitOrder(ctx context.Context, req ucEntity.InitOrderRequest) (ucEntity.InitOrderResponse, error)
}

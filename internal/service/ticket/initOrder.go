package ticket

import (
	"context"
	"errors"
	cacheTicket "go-projects/hexagonal-example/internal/adapter/outbound/cache/ticket"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
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
	// key reservasi per (user, event); cek reservasi untuk event ini saja.
	initOrder, err := s.Cache.Ticket.GetInitOrder(ctx, entity.CacheInitOrderRequest{UserID: userId, EventID: req.EventID})
	if err == nil {
		log.Info("reservation already exists, returning cached order",
			zap.String("cached_date", initOrder.Date))
		return ucEntity.InitOrderResponse{
			Date:     initOrder.Date,
			EventID:  initOrder.EventID,
			Quantity: initOrder.Quantity,
		}, nil
	}

	// parse tanggal & validasi event (sekaligus cek eksistensi + rentang tanggal)
	parsedEventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		log.Warn("invalid event date format", zap.String("date", req.Date), zap.Error(err))
		return response, err
	}

	event, err := s.Repository.Event.GetOneById(ctx, orm, req.ToObEvent(parsedEventDate))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("event not found")
			return response, errors.New("event not found")
		}
		log.Error("event validation failed", zap.Error(err))
		return response, err
	}

	// gate formulir pendaftaran (§5.3): kalau event punya form kustom dan user
	// belum mendaftar, pemesanan diblokir sampai formulir diisi. Cek registrasi
	// dulu; hanya kalau belum terdaftar baru cek apakah event memang butuh form.
	registered, regErr := s.Repository.UserRegistration.ExistsByUserEvent(ctx, orm, userId, req.EventID)
	if regErr != nil {
		log.Error("failed to check registration", zap.Error(regErr))
		return response, regErr
	}
	if !registered {
		fields, ferr := s.Repository.Event.GetFormFieldsByEvent(ctx, orm, req.EventID)
		if ferr != nil {
			log.Error("failed to check event form", zap.Error(ferr))
			return response, ferr
		}
		if len(fields) > 0 {
			log.Warn("registration required before init order")
			return response, errors.New("isi formulir pendaftaran terlebih dahulu")
		}
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

	// caching initialization data. harga diambil dari event (server-side),
	// bukan dari request, supaya tidak bisa dimanipulasi client.
	err = s.Cache.Ticket.SetInitOrder(ctx, req.ToObSetCache(userId, int64(event.Price)))
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

package ticket

import (
	"context"
	"strconv"
	"time"

	obEntity "go-projects/hexagonal-example/internal/adapter/outbound/entity"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

	"github.com/google/uuid"
)

func (s service) Purchase(ctx context.Context, req ucEntity.PurchaseRequest) (ucEntity.PurchaseResponse, error) {
	var (
		orm       = s.repository.DB
		response  ucEntity.PurchaseResponse
		userIdCtx = ctx.Value("x-user-id").(string)
		userId, _ = strconv.ParseInt(userIdCtx, 10, 64)
	)

	// ambil payload init order dari cache (di-set saat InitOrder)
	// event id diambil dari request; reservasi di-cache per (user, event).
	cachedInitOrder, err := s.Cache.Ticket.GetInitOrder(ctx, obEntity.CacheInitOrderRequest{UserID: userId, EventID: req.EventID})
	if err != nil {
		return response, err
	}

	parsedTime, err := time.Parse("2006-01-02", cachedInitOrder.Date)
	if err != nil {
		return response, err
	}

	// validasi event sekaligus ambil harga. GetOneById juga mengecek tanggal
	// pilihan user masih dalam rentang event.
	event, err := s.Repository.Event.GetOneById(ctx, orm, obEntity.Event{
		ID:        cachedInitOrder.EventID,
		StartDate: parsedTime,
	})
	if err != nil {
		return response, err
	}

	// masa berlaku tiket: pakai end date kalau event multi-hari, selain itu start date.
	validUntil := event.StartDate
	if !event.EndDate.IsZero() {
		validUntil = event.EndDate
	}

	// if event are free, then the transaction is automatically successful and tickets are issued immediately.
	status := obEntity.TransactionStatusPending
	if cachedInitOrder.Price == 0 {
		status = obEntity.TransactionStatusSuccessful
	}

	trx := orm.Begin()
	defer func() {
		if r := recover(); r != nil {
			trx.Rollback()
		}
	}()

	// tax/admin_fee/amount_deduction/promo dibiarkan nol/null; diisi setelah
	// call payment gateway tersedia.
	transaction := obEntity.Transaction{
		TxID:     uuid.NewString(),
		UserID:   userId,
		EventID:  event.ID,
		AuthorID: 1,
		Status:   status,
		Amount:   event.Price * float64(cachedInitOrder.Quantity),
	}
	if err = s.Repository.Transaction.Create(ctx, trx, transaction); err != nil {
		trx.Rollback()
		return response, err
	}

	// Free event, issue tickets immediately. Paid event, wait for payment gateway callback to issue tickets.
	if cachedInitOrder.Price == 0 {
		// decrement remaining quota in the database (atomic, locked di level baris).
		// satu UPDATE untuk seluruh quantity; RowsAffected 0 -> sold out -> rollback.
		if err = s.Repository.Event.DecrementQuota(ctx, trx, event.ID, cachedInitOrder.Quantity); err != nil {
			trx.Rollback()
			return response, err
		}

		for i := 0; i < int(cachedInitOrder.Quantity); i++ {
			ticket := obEntity.UserTicket{
				UserID:     userId,
				EventID:    event.ID,
				Code:       uuid.NewString(),
				Status:     obEntity.UserTicketStatusActive,
				ValidUntil: validUntil,
			}
			if err = s.Repository.UserTicket.Create(ctx, trx, ticket); err != nil {
				trx.Rollback()
				return response, err
			}
		}

	} else {
		// call payment gateway api (event berbayar) -- belum ada, dibiarkan
		// gateway request (log ke gateway_requests) -- menyusul bareng call gateway
		// TODO(callback): saat callback sukses, kurangi kuota DB via
		// Repository.Event.DecrementQuota(ctx, trx, event.ID, quantity) bareng
		// penerbitan tiket. Belum di sini karena tiket paid terbit di callback.
	}

	if err = trx.Commit().Error; err != nil {
		return response, err
	}

	// publish email notification to the user

	// Order sudah terbentuk; bersihkan cache init order.
	s.Cache.Ticket.ClearInitOrder(ctx, obEntity.CacheInitOrderRequest{UserID: userId, EventID: req.EventID})

	response = ucEntity.PurchaseResponse{Status: status}
	return response, nil
}

type IPurchase interface {
	Purchase(ctx context.Context, req ucEntity.PurchaseRequest) (ucEntity.PurchaseResponse, error)
}

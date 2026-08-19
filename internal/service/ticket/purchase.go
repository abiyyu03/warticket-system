package ticket

import (
	"context"
	"strconv"
	"time"

	obEntity "go-projects/hexagonal-example/internal/adapter/outbound/entity"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
	"go-projects/hexagonal-example/pkg/constants"

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
	cachedInitOrder, err := s.Cache.Ticket.GetInitOrder(ctx, obEntity.CacheInitOrderRequest{UserID: userId})
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

	trx := orm.Begin()
	defer func() {
		if r := recover(); r != nil {
			trx.Rollback()
		}
	}()

	// Buat transaksi berstatus PENDING. tax/admin_fee/amount_deduction/promo
	// dibiarkan nol/null; diisi setelah call payment gateway tersedia.
	transaction := obEntity.Transaction{
		TxID:     uuid.NewString(),
		UserID:   userId,
		EventID:  event.ID,
		AuthorID: 1,
		Status:   obEntity.TransactionStatusPending,
		Amount:   event.Price * float64(cachedInitOrder.Quantity),
	}
	if err = s.Repository.Transaction.Create(ctx, trx, transaction); err != nil {
		trx.Rollback()
		return response, err
	}

	// Terbitkan satu tiket per slot yang dibeli.
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

	// if the event was free
	if cachedInitOrder.Price == 0 {
		if err = s.Repository.Transaction.UpdateStatus(ctx, trx, obEntity.Transaction{
			TxID:   transaction.TxID,
			Status: constants.TicketOrderSuccess,
		}); err != nil {
			trx.Rollback()
			return response, err
		}

		response = ucEntity.PurchaseResponse{
			Status: constants.TicketOrderSuccess,
		}
	}

	// publish email notification to the user

	// call payment gateway api -- belum ada, dibiarkan
	// gateway request (log ke gateway_requests) -- menyusul bareng call gateway

	if err = trx.Commit().Error; err != nil {
		return response, err
	}

	// Order sudah terbentuk; bersihkan cache init order.
	s.Cache.Ticket.ClearInitOrder(ctx, obEntity.CacheInitOrderRequest{UserID: userId})

	response = ucEntity.PurchaseResponse{
		Status: constants.TicketOrderPending,
	}

	return response, nil
}

type IPurchase interface {
	Purchase(ctx context.Context, req ucEntity.PurchaseRequest) (ucEntity.PurchaseResponse, error)
}

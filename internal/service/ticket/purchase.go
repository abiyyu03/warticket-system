package ticket

import (
	"context"
	"strconv"

	obEntity "go-projects/hexagonal-example/internal/adapter/outbound/entity"
	ucEvent "go-projects/hexagonal-example/internal/service/entity/event"
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

	// ambil payload init order dari cache; reservasi di-key oleh tx_id yang
	// sudah diterbitkan saat InitOrder.
	cachedInitOrder, err := s.Cache.Ticket.GetInitOrder(ctx, obEntity.CacheInitOrderRequest{TxID: req.TxID})
	if err != nil {
		return response, err
	}

	// ambil event (tanpa validasi tanggal; date sudah ditentukan server-side).
	event, err := s.Repository.Event.GetByID(ctx, orm, cachedInitOrder.EventID)
	if err != nil {
		return response, err
	}

	// event ber-form: email peserta + jawaban formulir (payload checkout) wajib &
	// divalidasi di sini. Event tanpa form melewati langkah ini.
	formFields, err := s.Repository.Event.GetFormFieldsByEvent(ctx, orm, event.ID)
	if err != nil {
		return response, err
	}
	registration := ucEvent.SubmitRegistrationRequest{
		UserID:  userId,
		EventID: event.ID,
		Email:   req.Email,
		Answers: req.Answers,
	}
	if len(formFields) > 0 {
		if err = registration.Validate(formFields); err != nil {
			return response, err
		}
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
		TxID:     cachedInitOrder.TxID,
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

	// registrasi peserta (create-or-update) bila event ber-form; disimpan dalam
	// transaksi yang sama dengan order.
	if len(formFields) > 0 {
		reg := registration.ToObEntity()
		if err = s.Repository.UserRegistration.Upsert(ctx, trx, &reg); err != nil {
			trx.Rollback()
			return response, err
		}
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

		// TODO(email): distribusikan tiket ke email pendaftaran (belum ada login).
		// Email tersimpan wajib di user_registrations; ambil via (user, event) lalu
		// kirim lewat mailer. Pengiriman ditunda — mekanisme email menyusul.

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

	// Order sudah terbentuk; bersihkan cache init order (key tx_id).
	s.Cache.Ticket.ClearInitOrder(ctx, obEntity.CacheInitOrderRequest{TxID: cachedInitOrder.TxID})

	response = ucEntity.PurchaseResponse{Status: status}
	return response, nil
}

type IPurchase interface {
	Purchase(ctx context.Context, req ucEntity.PurchaseRequest) (ucEntity.PurchaseResponse, error)
}

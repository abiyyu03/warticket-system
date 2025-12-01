package ticket

import (
	"context"
	cryptoRand "crypto/rand"
	"fmt"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
	"io"
	"strconv"
	"strings"
	"time"
)

func (s service) ClaimAndPurchase(ctx context.Context, req ucEntity.ClaimTicketRequest) (ucEntity.ClaimAndPurchaseTicketResponse, error) {
	var (
		orm       = s.repository.DB
		response  ucEntity.ClaimAndPurchaseTicketResponse
		userIdCtx = ctx.Value("x-user-id").(string)
		userId, _ = strconv.ParseInt(userIdCtx, 10, 64)
	)

	// check cache availability if already exist, invoke
	cachedInitOrder, err := s.Cache.Ticket.GetInitOrder(ctx, req.ToObGetCache(userId))
	if err != nil {
		return response, err
	}

	parsedTime, err := time.Parse("2006-01-02", cachedInitOrder.Date)
	if err != nil {
		return response, err
	}

	trx := orm.Begin()
	func() {
		if r := recover(); r != nil {
			trx.Rollback()
		}
	}()

	// check chair availability
	for _, chairCode := range cachedInitOrder.ChairCode {
		seatAvailable, err := s.Repository.Seat.GetOne(ctx, orm, req.ToObSeat(chairCode))
		if err != nil {
			trx.Rollback()
			return response, err
		}

		eventDetail, err := s.Repository.EventDetail.GetOneByEventIdAndType(ctx, orm, req.ToObEventDetail(cachedInitOrder.EventID, seatAvailable.Type))
		if err != nil {
			trx.Rollback()
			return response, err
		}

		ticketCode, err := s.generateSecureRandomString(24)
		if err != nil {
			trx.Rollback()
			return response, err
		}

		err = s.Repository.UserEventTicket.Create(ctx, orm, req.ToObTicket(
			userId,
			seatAvailable.ID,
			eventDetail.ID,
			strings.ToUpper(ticketCode),
			parsedTime,
		))
		if err != nil {
			trx.Rollback()
			return response, err
		}

		// update seat status
		err = s.Repository.Seat.UpdateStatus(ctx, orm, req.ToObSeatUnavailable(seatAvailable.ID))
		if err != nil {
			trx.Rollback()
			return response, err
		}

		// revoke cache
		s.Cache.Ticket.ClearInitOrder(ctx, req.ToObGetCache(userId))

	}
	trx.Commit()

	response = ucEntity.ClaimAndPurchaseTicketResponse{
		Status: "TICKET_ACQUIRED",
	}

	return response, nil
}

func (s service) generateSecureRandomString(length int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const maxIndex = len(letters)

	// Create a slice of bytes to hold the random indexes
	b := make([]byte, length)

	// Read random data from the crypto/rand source
	if _, err := io.ReadFull(cryptoRand.Reader, b); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	// Map the random bytes to the characters in the 'letters' string
	for i := 0; i < length; i++ {
		// Use the modulo operator (%) to map the random byte (0-255)
		// to an index within the 'letters' string (0-61).
		// This introduces a slight bias but is acceptable for most non-cryptographic keys.
		b[i] = letters[int(b[i])%maxIndex]
	}

	return string(b), nil
}

type IClaimAndPurchase interface {
	ClaimAndPurchase(ctx context.Context, req ucEntity.ClaimTicketRequest) (ucEntity.ClaimAndPurchaseTicketResponse, error)
}

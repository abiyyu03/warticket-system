package ticket

import (
	"context"
	cryptoRand "crypto/rand"
	"fmt"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
	"io"
	"math/rand/v2"
	"time"
)

func (s service) Claim(ctx context.Context, req ucEntity.ClaimTicketRequest) (ucEntity.ClaimTicketResponse, error) {
	var (
		orm      = s.repository.DB
		response ucEntity.ClaimTicketResponse
	)

	// check event available
	event, err := s.Repository.Event.GetOneById(ctx, orm, req.ToObEvent())
	if err != nil {
		return response, err
	}

	parsedTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return response, err
	}

	// check chaircode
	if req.ChairCode != nil {
		eventDetail, err := s.Repository.EventDetail.GetOneByEventIdAndType(ctx, orm, req.ToObEventDetail("VIP"))
		if err != nil {
			return response, err
		}

		// check chair availability
		seatAvailable, err := s.Repository.Seat.GetOne(ctx, orm, req.ToObSeat())
		if err != nil {
			return response, err
		}

		err = s.Repository.UserEventTicket.Create(ctx, orm, req.ToObTicket(
			1,
			seatAvailable.ID,
			eventDetail.ID,
			eventDetail.TicketCode,
			parsedTime,
		))
		if err != nil {
			return response, err
		}

		response = ucEntity.ClaimTicketResponse{
			Status: "Ticket VIP Found",
		}
	} else {
		// random available seat
		seats, err := s.Repository.Seat.GetSeatsByLocation(ctx, orm, req.ToObSeatLocation(event.LocationID))
		if err != nil {
			return response, err
		}

		randomNumber := rand.IntN(len(seats)) + 1

		// WARNING : POSSIBILITY RACE CONDITION
		eventDetail, err := s.Repository.EventDetail.GetOneByEventIdAndType(ctx, orm, req.ToObEventDetail("REGULAR"))
		if err != nil {
			return response, err
		}

		// get chair with randomize method

		err = s.Repository.UserEventTicket.Create(ctx, orm, req.ToObTicket(
			1,
			seatAvailable.ID,
			eventDetail.ID,
			eventDetail.TicketCode,
			parsedTime,
		))
		if err != nil {
			return response, err
		}

		response = ucEntity.ClaimTicketResponse{
			Status: "Ticket Found",
		}
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

type IClaim interface {
	Claim(ctx context.Context, req ucEntity.ClaimTicketRequest) (ucEntity.ClaimTicketResponse, error)
}

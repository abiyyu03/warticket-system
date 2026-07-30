package ticket

import (
	obEntity "go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"time"
)

type (
	RedeemRequest struct {
		TicketCode string `json:"ticket_code"`
	}
)

func (r RedeemRequest) ToObEntityRedeem(redeemTimestamp *time.Time) obEntity.UserEventTicket {
	return obEntity.UserEventTicket{
		TicketCode: r.TicketCode,
		IsRedeem:   true,
		RedeemedAt: redeemTimestamp,
	}
}

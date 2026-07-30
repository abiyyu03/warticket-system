package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

type (
	RedeemRequest struct {
		TicketCode string `json:"ticket_code"`
	}
)

func (r RedeemRequest) ToUcEntity() ucEntity.RedeemRequest {
	return ucEntity.RedeemRequest{
		TicketCode: r.TicketCode,
	}
}

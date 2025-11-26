package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

type (
	ClaimTicketRequest struct {
		ChairCode *int   `json:"chair_code"` // book for vip only
		Date      string `json:"date"`
		EventID   int64  `json:"event_id"`
	}
)

func (r ClaimTicketRequest) ToUcEntity() ucEntity.ClaimTicketRequest {
	return ucEntity.ClaimTicketRequest{}
}

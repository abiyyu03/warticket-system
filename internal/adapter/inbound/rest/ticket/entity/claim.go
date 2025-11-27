package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

type (
	ClaimTicketRequest struct {
		ChairCode *string `json:"chair_code"` // book for vip only
		Date      string  `json:"date"`
		EventID   int64   `json:"event_id"`
	}
)

func (r ClaimTicketRequest) ToUcEntity() ucEntity.ClaimTicketRequest {
	result := ucEntity.ClaimTicketRequest{
		Date:    r.Date,
		EventID: r.EventID,
	}

	if r.ChairCode != nil {
		result.ChairCode = r.ChairCode
	}

	return result
}

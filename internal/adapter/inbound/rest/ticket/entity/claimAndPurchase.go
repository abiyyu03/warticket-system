package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

type (
	ClaimTicketRequest struct{}
)

func (r ClaimTicketRequest) ToUcEntity() ucEntity.ClaimTicketRequest {
	return ucEntity.ClaimTicketRequest{}
}

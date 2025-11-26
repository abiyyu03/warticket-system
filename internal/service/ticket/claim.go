package ticket

import (
	"context"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
)

func (s service) Claim(ctx context.Context, req ucEntity.ClaimTicketRequest) error {

	return nil
}

type IClaim interface {
	Claim(ctx context.Context, req ucEntity.ClaimTicketRequest) error
}

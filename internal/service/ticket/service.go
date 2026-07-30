package ticket

import (
	"go-projects/hexagonal-example/internal/adapter/outbound"
	"go-projects/hexagonal-example/pkg"
)

type TicketService interface {
	IClaimAndPurchase
	IRedeem
	IInitOrder
}

type service struct {
	outbound.Outbound
	cache      *pkg.Redis
	repository *pkg.SQL
}

func New(
	outbound outbound.Outbound,
	pkg pkg.Package,
) TicketService {
	return &service{
		Outbound:   outbound,
		cache:      pkg.Cache,
		repository: pkg.DB,
	}
}

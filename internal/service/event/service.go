package event

import (
	"go-projects/hexagonal-example/internal/adapter/outbound"
	"go-projects/hexagonal-example/pkg"
)

type EventService interface {
	ICreateEvent
	IGetListEvent
}

type service struct {
	outbound.Outbound
	cache      *pkg.Redis
	repository *pkg.SQL
}

func New(
	outbound outbound.Outbound,
	pkg pkg.Package,
) EventService {
	return &service{
		Outbound:   outbound,
		cache:      pkg.Cache,
		repository: pkg.DB,
	}
}

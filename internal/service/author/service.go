package author

import (
	"go-projects/hexagonal-example/internal/adapter/outbound"
	"go-projects/hexagonal-example/pkg"
	"go-projects/hexagonal-example/pkg/token"
)

type AuthorService interface {
	IRegister
	ILogin
	IRefresh
}

type service struct {
	outbound.Outbound
	repository *pkg.SQL
	token      *token.Manager
}

func New(
	outbound outbound.Outbound,
	pkg pkg.Package,
) AuthorService {
	return &service{
		Outbound:   outbound,
		repository: pkg.DB,
		token:      token.New(),
	}
}

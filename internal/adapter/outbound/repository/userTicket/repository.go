package userTicket

import "go-projects/hexagonal-example/pkg"

type Repository interface {
	ICreate
	IRedeemTicket
	IGetOneByCode
}

type userTicket struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Repository {
	return &userTicket{
		Package: pkg,
	}
}

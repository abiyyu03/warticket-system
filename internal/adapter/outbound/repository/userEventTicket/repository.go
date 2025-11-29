package userEventTicket

import "go-projects/hexagonal-example/pkg"

type Repository interface {
	ICreate
	IUpdateByCode
}

type userEventTicket struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Repository {
	return &userEventTicket{
		Package: pkg,
	}
}

package userRegistration

import "go-projects/hexagonal-example/pkg"

type Repository interface {
	ICreate
	IUpsert
}

type userRegistration struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Repository {
	return &userRegistration{
		Package: pkg,
	}
}

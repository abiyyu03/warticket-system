package userEvent

import "go-projects/hexagonal-example/pkg"

type Repository interface {
	ICreate
}

type user struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Repository {
	return &user{
		Package: pkg,
	}
}

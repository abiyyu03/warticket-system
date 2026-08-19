package event

import "go-projects/hexagonal-example/pkg"

type Repository interface {
	IGetOneById
	IGetOneByCode
	ICreate
	IGetAll
}

type event struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Repository {
	return &event{
		Package: pkg,
	}
}

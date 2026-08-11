package transaction

import "go-projects/hexagonal-example/pkg"

type Repository interface {
	ICreate
}

type transaction struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Repository {
	return &transaction{
		Package: pkg,
	}
}

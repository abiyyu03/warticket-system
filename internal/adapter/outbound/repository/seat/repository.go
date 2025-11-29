package seat

import "go-projects/hexagonal-example/pkg"

type Repository interface {
	IGetOne
	IUpdateByCode
	IUpdateStatus
	IGetSeatsByLocation
}

type seat struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Repository {
	return &seat{
		Package: pkg,
	}
}

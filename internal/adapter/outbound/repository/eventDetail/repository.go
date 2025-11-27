package eventDetail

import "go-projects/hexagonal-example/pkg"

type Repository interface {
	IGetOneByEventIdAndType
}

type eventDetail struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Repository {
	return &eventDetail{
		Package: pkg,
	}
}

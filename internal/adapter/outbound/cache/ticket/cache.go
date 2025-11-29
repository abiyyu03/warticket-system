package ticket

import "go-projects/hexagonal-example/pkg"

var (
	key = "tickets:initOrder:%d"
)

type Cache interface {
	IGetInitOrder
	ISetInitOrder
	IClearInitOrder
}

type ticketCache struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Cache {
	return &ticketCache{
		Package: pkg,
	}
}

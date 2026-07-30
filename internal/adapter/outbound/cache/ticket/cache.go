package ticket

import "go-projects/hexagonal-example/pkg"

var (
	key          = "tickets:order:%d"
	keyInit      = "tickets:order:init:%d"
	keyDecrQuota = "tickets:event:%d"
)

type Cache interface {
	IGetInitOrder
	ISetInitOrder
	IClearInitOrder
	IDecrTicketQuota
	ISetTicketQuota
}

type ticketCache struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Cache {
	return &ticketCache{
		Package: pkg,
	}
}

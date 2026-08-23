package ticket

import "go-projects/hexagonal-example/pkg"

var (
	// reservasi init order di-key oleh tx_id (diterbitkan saat init order).
	key          = "tickets:order:%s"
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

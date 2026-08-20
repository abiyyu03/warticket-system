package event

import (
	"errors"

	"go-projects/hexagonal-example/pkg"
)

var ErrInsufficientQuota = errors.New("insufficient event quota")

type Repository interface {
	IGetOneById
	IGetOneByCode
	ICreate
	IGetAll
	IDecrementQuota
}

type event struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Repository {
	return &event{
		Package: pkg,
	}
}

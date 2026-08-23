package userAuthor

import "go-projects/hexagonal-example/pkg"

type Repository interface {
	ICreate
	IGetByEmail
}

type userAuthor struct {
	Package pkg.Package
}

func New(pkg pkg.Package) Repository {
	return &userAuthor{
		Package: pkg,
	}
}

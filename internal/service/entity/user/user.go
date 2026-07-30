package user

import "go-projects/hexagonal-example/internal/adapter/outbound/entity"

type User struct {
	ID    int64
	Email string
	Name  string
}

func (r User) ToOutbondUser() entity.User {
	return entity.User{
		Email: r.Email,
		Name:  r.Name,
	}
}

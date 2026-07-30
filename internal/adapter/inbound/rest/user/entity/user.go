package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/user"

type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (r CreateUserRequest) ToServiceUser() ucEntity.User {
	return ucEntity.User{
		Email: r.Email,
		Name:  r.Name,
	}
}

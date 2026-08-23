package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/author"

type (
	RegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}
)

func (r RegisterRequest) ToUcEntity() ucEntity.RegisterRequest {
	return ucEntity.RegisterRequest{Name: r.Name, Email: r.Email, Password: r.Password}
}

func (r LoginRequest) ToUcEntity() ucEntity.LoginRequest {
	return ucEntity.LoginRequest{Email: r.Email, Password: r.Password}
}

func (r RefreshRequest) ToUcEntity() ucEntity.RefreshRequest {
	return ucEntity.RefreshRequest{RefreshToken: r.RefreshToken}
}

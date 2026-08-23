package author

import (
	"context"
	"errors"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/author"
	"go-projects/hexagonal-example/pkg/token"

	"golang.org/x/crypto/bcrypt"
)

// errInvalidCredential sengaja generik supaya tidak membocorkan mana yang salah.
var errInvalidCredential = errors.New("email atau password salah")

// Login memverifikasi kredensial author dan menerbitkan pasangan token (role author).
func (s service) Login(ctx context.Context, req ucEntity.LoginRequest) (ucEntity.AuthResponse, error) {
	var response ucEntity.AuthResponse

	author, err := s.Repository.UserAuthor.GetByEmail(ctx, s.repository.DB, req.Email)
	if err != nil {
		return response, errInvalidCredential
	}

	if bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(req.Password)) != nil {
		return response, errInvalidCredential
	}

	access, refresh, expiresIn, err := s.token.GeneratePair(author.ID, token.RoleAuthor)
	if err != nil {
		return response, err
	}

	return ucEntity.AuthResponse{
		TokenType:    "Bearer",
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		Role:         token.RoleAuthor,
	}, nil
}

type ILogin interface {
	Login(ctx context.Context, req ucEntity.LoginRequest) (ucEntity.AuthResponse, error)
}

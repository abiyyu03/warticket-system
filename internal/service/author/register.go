package author

import (
	"context"
	"errors"
	obEntity "go-projects/hexagonal-example/internal/adapter/outbound/entity"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/author"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register membuat akun author baru; password di-hash bcrypt.
func (s service) Register(ctx context.Context, req ucEntity.RegisterRequest) error {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" || req.Password == "" {
		return errors.New("name, email, dan password wajib diisi")
	}

	// cegah email ganda (backstop: unique constraint).
	if _, err := s.Repository.UserAuthor.GetByEmail(ctx, s.repository.DB, req.Email); err == nil {
		return errors.New("email sudah terdaftar")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	author := obEntity.UserAuthor{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
	}
	return s.Repository.UserAuthor.Create(ctx, s.repository.DB, &author)
}

type IRegister interface {
	Register(ctx context.Context, req ucEntity.RegisterRequest) error
}

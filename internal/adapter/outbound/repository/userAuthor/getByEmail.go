package userAuthor

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

// GetByEmail mengambil author berdasarkan email; gorm.ErrRecordNotFound bila tidak ada.
func (r userAuthor) GetByEmail(ctx context.Context, orm *gorm.DB, email string) (entity.UserAuthor, error) {
	var author entity.UserAuthor
	err := orm.WithContext(ctx).Where("email = ?", email).First(&author).Error
	return author, err
}

type IGetByEmail interface {
	GetByEmail(ctx context.Context, orm *gorm.DB, email string) (entity.UserAuthor, error)
}

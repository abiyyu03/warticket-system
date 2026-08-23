package userAuthor

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r userAuthor) Create(ctx context.Context, orm *gorm.DB, author *entity.UserAuthor) error {
	return orm.WithContext(ctx).Create(author).Error
}

type ICreate interface {
	Create(ctx context.Context, orm *gorm.DB, author *entity.UserAuthor) error
}

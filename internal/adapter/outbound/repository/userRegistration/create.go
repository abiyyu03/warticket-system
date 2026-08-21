package userRegistration

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r userRegistration) Create(ctx context.Context, orm *gorm.DB, reg *entity.UserRegistration) error {
	return orm.WithContext(ctx).Create(reg).Error
}

type ICreate interface {
	Create(ctx context.Context, orm *gorm.DB, reg *entity.UserRegistration) error
}

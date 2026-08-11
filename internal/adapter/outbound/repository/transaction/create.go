package transaction

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r transaction) Create(ctx context.Context, orm *gorm.DB, trx entity.Transaction) error {
	err := orm.WithContext(ctx).Create(&trx).Error
	if err != nil {
		return err
	}

	return nil
}

type ICreate interface {
	Create(ctx context.Context, orm *gorm.DB, trx entity.Transaction) error
}

package event

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r event) DecrementQuota(ctx context.Context, orm *gorm.DB, eventID, qty int64) error {
	res := orm.WithContext(ctx).
		Model(&entity.Event{}).
		Where("id = ? AND quota_remaining >= ?", eventID, qty).
		Update("quota_remaining", gorm.Expr("quota_remaining - ?", qty))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrInsufficientQuota
	}
	return nil
}

type IDecrementQuota interface {
	DecrementQuota(ctx context.Context, orm *gorm.DB, eventID, qty int64) error
}
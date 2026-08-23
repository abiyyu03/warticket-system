package userRegistration

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

// ExistsByUserEvent memeriksa apakah seorang user sudah mendaftar pada event.
// Pakai COUNT (bukan First) supaya "belum terdaftar" — kondisi normal di gate
// init-order — tidak memicu log ErrRecordNotFound dari GORM.
func (r userRegistration) ExistsByUserEvent(ctx context.Context, orm *gorm.DB, userID, eventID int64) (bool, error) {
	var count int64
	err := orm.WithContext(ctx).
		Model(&entity.UserRegistration{}).
		Where("user_id = ? AND event_id = ?", userID, eventID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

type IExistsByUserEvent interface {
	ExistsByUserEvent(ctx context.Context, orm *gorm.DB, userID, eventID int64) (bool, error)
}

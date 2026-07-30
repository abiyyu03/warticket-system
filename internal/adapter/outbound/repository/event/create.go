package event

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

// Create menyimpan event. event dilewatkan sebagai pointer supaya ID hasil
// auto-increment terisi balik dan bisa dipakai memberi key counter kuota di Redis.
func (r event) Create(ctx context.Context, orm *gorm.DB, event *entity.Event) error {
	err := orm.WithContext(ctx).Create(event).Error
	if err != nil {
		return err
	}

	return nil
}

type ICreate interface {
	Create(ctx context.Context, orm *gorm.DB, event *entity.Event) error
}

package event

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

// CreateFormFields menyimpan definisi field formulir sebuah event secara bulk.
// Dipanggil dalam transaksi yang sama dengan pembuatan event.
func (r event) CreateFormFields(ctx context.Context, orm *gorm.DB, fields []entity.EventFormField) error {
	if len(fields) == 0 {
		return nil
	}
	return orm.WithContext(ctx).Create(&fields).Error
}

type ICreateFormFields interface {
	CreateFormFields(ctx context.Context, orm *gorm.DB, fields []entity.EventFormField) error
}

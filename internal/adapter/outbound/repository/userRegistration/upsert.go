package userRegistration

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Upsert menyimpan registrasi; bila (user_id, event_id) sudah ada, email &
// answers diperbarui. Memanfaatkan unique uq_user_registrations_user_event.
func (r userRegistration) Upsert(ctx context.Context, orm *gorm.DB, reg *entity.UserRegistration) error {
	return orm.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "event_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"email", "answers", "updated_at"}),
		}).
		Create(reg).Error
}

type IUpsert interface {
	Upsert(ctx context.Context, orm *gorm.DB, reg *entity.UserRegistration) error
}

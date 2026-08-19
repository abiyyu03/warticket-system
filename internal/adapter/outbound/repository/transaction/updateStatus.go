package transaction

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r transaction) UpdateStatus(ctx context.Context, orm *gorm.DB, trx entity.Transaction) error {
	err := orm.WithContext(ctx).Where("tx_id = ?", trx.TxID).Updates(map[string]interface{}{"status": trx.Status}).Error
	if err != nil {
		return err
	}

	return nil
}

type IUpdateStatus interface {
	UpdateStatus(ctx context.Context, orm *gorm.DB, trx entity.Transaction) error
}

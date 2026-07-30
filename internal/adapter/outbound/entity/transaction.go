package entity

import (
	"encoding/json"
)

// Status transaksi. Nilainya dibatasi CHECK constraint ck_transactions_status,
// jadi konstanta di sini harus ikut diubah kalau migration-nya berubah.
const (
	TransactionStatusPending  = "PENDING"
	TransactionStatusPaid     = "PAID"
	TransactionStatusFailed   = "FAILED"
	TransactionStatusExpired  = "EXPIRED"
	TransactionStatusRefunded = "REFUNDED"
)

type Transaction struct {
	BaseModel
	ID      int64   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TrxID   string  `gorm:"column:trx_id;not null;uniqueIndex" json:"trx_id"`
	UserID  int64   `gorm:"column:user_id;not null;index" json:"user_id"`
	EventID int64   `gorm:"column:event_id;not null;index" json:"event_id"`
	Amount  float64 `gorm:"column:amount;not null" json:"amount"`
	Status  string  `gorm:"column:status;not null;default:PENDING;index" json:"status"`
	User    User    `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Event   Event   `gorm:"foreignKey:EventID;references:ID" json:"event"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func (u Transaction) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u Transaction) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

package entity

import (
	"encoding/json"
	"time"
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
	ID       int64   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TxID     string  `gorm:"column:tx_id;not null;uniqueIndex" json:"tx_id"`
	UserID   int64   `gorm:"column:user_id;not null;index" json:"user_id"`
	EventID  int64   `gorm:"column:event_id;not null;index" json:"event_id"`
	AuthorID int64   `gorm:"column:author_id;not null;index" json:"author_id"`
	Status   string  `gorm:"column:status;not null;default:PENDING;index" json:"status"`
	Amount   float64 `gorm:"column:amount;not null" json:"amount"`
	// Nullable: kosong kalau tidak ada potongan.
	AmountDeduction *float64 `gorm:"column:amount_deduction" json:"amount_deduction,omitempty"`
	// Nullable: terisi kalau transaksi memakai promo.
	PromoID  *int64  `gorm:"column:promo_id" json:"promo_id,omitempty"`
	Tax      float64 `gorm:"column:tax;not null;default:0" json:"tax"`
	AdminFee float64 `gorm:"column:admin_fee;not null;default:0" json:"admin_fee"`
	// Nullable: terisi saat pembayaran berhasil.
	PaymentAt *time.Time `gorm:"column:payment_at" json:"payment_at,omitempty"`

	User   User  `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Event  Event `gorm:"foreignKey:EventID;references:ID" json:"event"`
	Author User  `gorm:"foreignKey:AuthorID;references:ID" json:"author"`
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

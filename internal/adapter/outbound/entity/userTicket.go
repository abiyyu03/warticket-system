package entity

import (
	"encoding/json"
	"time"
)

// Status tiket. Nilainya dibatasi CHECK constraint ck_user_tickets_status.
const (
	UserTicketStatusActive    = "ACTIVE"
	UserTicketStatusRedeemed  = "REDEEMED"
	UserTicketStatusExpired   = "EXPIRED"
	UserTicketStatusCancelled = "CANCELLED"
)

type UserTicket struct {
	BaseModel
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"column:user_id;not null;index" json:"user_id"`
	EventID    int64     `gorm:"column:event_id;not null;index" json:"event_id"`
	Code       string    `gorm:"column:code;not null;uniqueIndex" json:"code"`
	Status     string    `gorm:"column:status;not null;default:ACTIVE;index" json:"status"`
	ValidUntil time.Time `gorm:"column:valid_until;not null" json:"valid_until"`
	User       User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Event      Event     `gorm:"foreignKey:EventID;references:ID" json:"event"`
}

func (UserTicket) TableName() string {
	return "user_tickets"
}

func (u UserTicket) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u UserTicket) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

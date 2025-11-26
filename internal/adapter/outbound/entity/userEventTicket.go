package entity

import (
	"encoding/json"
	"time"
)

type UserEventTicket struct {
	BaseModel
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        int64      `gorm:"column:user_id;not null;index" json:"user_id"`
	SeatID        int64      `gorm:"column:seat_id;not null;index" json:"seat_id"`
	EventDetailID int64      `gorm:"column:event_detail_id;not null" json:"event_detail_id"`
	TicketCode    string     `gorm:"column:ticket_code;not null" json:"ticket_code"`
	IsRedeem      bool       `gorm:"column:is_redeem;not null;default:false" json:"is_redeem"`
	RedeemedAt    *time.Time `gorm:"column:redeemed_at;not null" json:"redeemed_at"`
}

func (UserEventTicket) TableName() string {
	return "user_event_tickets"
}

func (u UserEventTicket) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u UserEventTicket) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

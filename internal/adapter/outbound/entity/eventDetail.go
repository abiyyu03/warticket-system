package entity

import (
	"encoding/json"
)

type EventDetail struct {
	BaseModel
	ID            int64   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	EventID       int64   `gorm:"column:event_id;not null" json:"event_id"`
	TicketCode    string  `gorm:"column:ticket_code;not null;uniqueIndex" json:"ticket_code"`
	TicketType    string  `gorm:"column:ticket_type;not null;default:REGULAR" json:"ticket_type"` // VIP or REGULAR
	Price         float64 `gorm:"column:price;not null;check:price >= 0" json:"price"`
	AvailableSlot int64   `gorm:"column:available_slot;not null;check:available_slot > 0" json:"available_slot"`
	Status        string  `gorm:"column:status;not null;default:'active';index" json:"status"`
	Event         Event   `gorm:"foreignKey:EventID;references:ID" json:"event"`
}

func (EventDetail) TableName() string {
	return "event_details"
}

func (u EventDetail) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u EventDetail) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

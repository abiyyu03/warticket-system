package entity

import (
	"encoding/json"
	"time"
)

type Event struct {
	BaseModel
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Description string    `gorm:"column:description;type:text" json:"description,omitempty"`
	ImageFile   string    `gorm:"column:image_file" json:"image_file,omitempty"`
	Price       float64   `gorm:"column:price;not null" json:"price"`
	Quota       int64     `gorm:"column:quota;not null" json:"quota"`
	QuotaRemaining int64  `gorm:"column:quota_remaining;not null" json:"quota_remaining"`
	StartDate   time.Time `gorm:"column:start_date;not null;index" json:"start_date"`
	EndDate     time.Time `gorm:"column:end_date" json:"end_date"`
}

func (Event) TableName() string {
	return "events"
}

func (u Event) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u Event) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

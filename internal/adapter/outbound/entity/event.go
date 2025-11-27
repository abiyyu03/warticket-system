package entity

import (
	"encoding/json"
	"time"
)

type Event struct {
	BaseModel
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"column:name;not null;index" json:"name"`
	Cover       string     `gorm:"column:cover;not null;index" json:"cover"`
	LocationID  int64      `gorm:"column:location_id;not null" json:"location_id"`
	Description string     `gorm:"column:description;type:text" json:"description,omitempty"`
	Status      string     `gorm:"column:status;not null;default:'active';index" json:"status"`
	StartDate   time.Time  `gorm:"column:start_date;not null;index" json:"start_date"`
	EndDate     *time.Time `gorm:"column:end_date" json:"end_date"`
	Location    Location   `gorm:"foreignKey:LocationID;references:ID" json:"location"`
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

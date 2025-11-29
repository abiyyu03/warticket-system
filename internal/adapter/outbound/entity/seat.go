package entity

import (
	"encoding/json"
)

type (
	Seat struct {
		BaseModel
		ID         int64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
		Code       string   `gorm:"column:code;not null;uniqueIndex" json:"code"`
		Row        int      `gorm:"column:row;not null" json:"row"`
		Col        int      `gorm:"column:col;not null" json:"col"`
		Available  bool     `gorm:"column:available;not null;default:true" json:"available"`
		Type       string   `gorm:"column:type;not null;default:REGULAR" json:"type"` // VIP or REGULAR
		LocationID int64    `gorm:"column:location_id;not null" json:"location_id"`
		Location   Location `gorm:"foreignKey:LocationID;references:ID" json:"location"`
	}
)

func (Seat) TableName() string {
	return "seats"
}

func (u Seat) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u Seat) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

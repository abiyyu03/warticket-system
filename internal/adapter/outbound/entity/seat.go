package entity

import (
	"encoding/json"
)

type (
	Seat struct {
		BaseModel
		ID        int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
		Code      string `gorm:"column:code;not null;uniqueIndex" json:"code"`
		Available string `gorm:"column:available;not null;default:true" json:"available"`
		Type      string `gorm:"column:type;not null;default:REGULAR" json:"type"` // VIP or REGULAR
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

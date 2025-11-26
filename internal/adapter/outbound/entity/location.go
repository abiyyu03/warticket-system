package entity

import (
	"encoding/json"
)

type (
	LocCoordinate struct {
		Lat float64 `gorm:"column:latitude;type:decimal(10,8);not null;index" json:"lat"`
		Lng float64 `gorm:"column:longitude;type:decimal(11,8);not null;index" json:"lng"`
	}

	Location struct {
		BaseModel
		ID         int64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
		Name       string        `gorm:"column:name;not null;index" json:"name"`
		Coordinate LocCoordinate `gorm:"embedded" json:"coordinate"`
	}
)

func (Location) TableName() string {
	return "locations"
}

func (u Location) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u Location) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

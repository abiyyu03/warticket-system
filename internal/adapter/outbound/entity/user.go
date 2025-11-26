package entity

import (
	"encoding/json"
)

type User struct {
	BaseModel
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Email    string `gorm:"column:email;uniqueIndex;not null" json:"email"`
	Name     string `gorm:"column:name;not null" json:"name"`
	Password string `gorm:"column:password;not null" json:"password"`
}

func (User) TableName() string {
	return "users"
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

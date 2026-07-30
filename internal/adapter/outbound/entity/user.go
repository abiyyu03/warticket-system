package entity

import (
	"encoding/json"
)

type User struct {
	BaseModel
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name;not null" json:"name"`
	Email    string `gorm:"column:email;not null;uniqueIndex" json:"email"`
	Password string `gorm:"column:password;not null" json:"password"`
	IsBanned bool   `gorm:"column:is_banned;not null;default:false" json:"is_banned"`
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

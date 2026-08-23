package entity

import "encoding/json"

type UserAuthor struct {
	BaseModel
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name;not null" json:"name"`
	Email    string `gorm:"column:email;not null;uniqueIndex" json:"email"`
	Password string `gorm:"column:password;not null" json:"-"` // hash bcrypt, tidak diekspos
}

func (UserAuthor) TableName() string {
	return "user_authors"
}

func (u UserAuthor) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u UserAuthor) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}

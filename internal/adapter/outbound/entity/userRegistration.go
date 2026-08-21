package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// AnswerMap dipetakan ke kolom JSONB: field_id (string) -> daftar nilai jawaban.
type AnswerMap map[string][]string

func (a AnswerMap) Value() (driver.Value, error) {
	if a == nil {
		return json.Marshal(AnswerMap{})
	}
	return json.Marshal(a)
}

func (a *AnswerMap) Scan(src any) error {
	if src == nil {
		*a = AnswerMap{}
		return nil
	}

	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("AnswerMap: tipe tidak didukung %T", src)
	}

	return json.Unmarshal(b, a)
}

type UserRegistration struct {
	BaseModel
	ID      int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID  int64     `gorm:"column:user_id;not null;index" json:"user_id"`
	EventID int64     `gorm:"column:event_id;not null;index" json:"event_id"`
	Answers AnswerMap `gorm:"column:answers;type:jsonb" json:"answers"`
}

func (UserRegistration) TableName() string {
	return "user_registrations"
}

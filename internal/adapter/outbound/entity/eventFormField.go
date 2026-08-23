package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Tipe field formulir pendaftaran. Dibatasi CHECK ck_event_form_fields_type,
// jadi konstanta di sini harus sinkron dengan migration.
const (
	EventFieldTypeText     = "text"     // isian teks singkat
	EventFieldTypeSelect   = "select"   // pilihan tunggal
	EventFieldTypeCheckbox = "checkbox" // pilihan ganda
)

// StringList dipetakan ke kolom JSONB (array string). Dipakai untuk daftar opsi
// field pilihan. Implement driver.Valuer + sql.Scanner supaya bisa dibaca/tulis
// GORM tanpa dependency tambahan.
type StringList []string

func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringList) Scan(src any) error {
	if src == nil {
		*s = nil
		return nil
	}

	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("StringList: tipe tidak didukung %T", src)
	}

	return json.Unmarshal(b, s)
}

type EventFormField struct {
	BaseModel
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	EventID   int64      `gorm:"column:event_id;not null;index" json:"event_id"`
	Label     string     `gorm:"column:label;not null" json:"label"`
	FieldType string     `gorm:"column:field_type;not null" json:"field_type"`
	Required  bool       `gorm:"column:required;not null;default:false" json:"required"`
	Options   StringList `gorm:"column:options;type:jsonb" json:"options,omitempty"`
	Position  int        `gorm:"column:position;not null;default:0" json:"position"`
}

func (EventFormField) TableName() string {
	return "event_form_fields"
}

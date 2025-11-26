package entity

import (
	"database/sql"
	"time"
)

type (
	BaseModel struct {
		CreatedAt time.Time    `gorm:"autoCreateTime:milli" json:"created_at"`
		UpdatedAt time.Time    `gorm:"autoUpdateTime:milli" json:"updated_at"`
		DeletedAt sql.NullTime `gorm:"index" json:"deleted_at,omitempty"`
	}
)

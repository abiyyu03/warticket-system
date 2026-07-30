package entity

import (
	"time"
)

type (
	// BaseModel memuat kolom waktu yang dipakai semua tabel.
	// Skema tidak lagi memakai soft delete, jadi tidak ada DeletedAt di sini.
	BaseModel struct {
		CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	}
)

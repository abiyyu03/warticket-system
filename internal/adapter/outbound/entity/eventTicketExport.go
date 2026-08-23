package entity

import "time"

type EventTicketExport struct {
	ID         int64
	Code       string
	UserID     int64
	Email      string
	Status     string
	ValidUntil time.Time
	CreatedAt  time.Time
}

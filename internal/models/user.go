package models

import "time"

// User represents the data structure for a user in the system.
type User struct {
	ID           int64
	Name         string
	Email        string
	TelegramID   int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastNotified time.Time
}

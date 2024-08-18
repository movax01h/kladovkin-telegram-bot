package models

import "time"

// User represents the user subscribing to notifications.
type User struct {
	ID           int64     `json:"id"`
	TelegramID   int64     `json:"telegram_id"`
	UserName     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	LastNotified time.Time `json:"last_notified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

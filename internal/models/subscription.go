package models

import "time"

// Subscription represents a user's subscription to a unit.
type Subscription struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	City      string    `json:"city"`
	Storage   string    `json:"storage"`
	UnitSize  string    `json:"unit_size"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

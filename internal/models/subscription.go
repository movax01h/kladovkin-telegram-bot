package models

import "time"

// Subscription represents the data structure for a subscription in the system.
type Subscription struct {
	ID        int64
	UserID    int64
	UnitID    int64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

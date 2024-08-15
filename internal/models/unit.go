package models

import "time"

// Unit represents the data structure for a unit in the system.
type Unit struct {
	ID          int64
	Name        string
	City        string
	Size        string
	Dimension   string
	Price       float64
	Available   bool
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

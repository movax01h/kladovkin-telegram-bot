package models

import "time"

// Unit represents a storage unit.
type Unit struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Size        string    `json:"size"`
	Dimension   string    `json:"dimension"`
	Price       float64   `json:"price"`
	Available   bool      `json:"available"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

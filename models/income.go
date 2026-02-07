package models

import "time"

type Income struct {
	ID        int64     `json:"id"`
	Year      int       `json:"year"`
	Month     int       `json:"month"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

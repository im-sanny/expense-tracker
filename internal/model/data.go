package model

import "time"

type Expense struct {
	ID     int       `json:"id" db:"id"`
	Date   time.Time `json:"date" db:"date"`
	Amount *int64    `json:"amount" db:"amount"`
	Note   *string   `json:"note" db:"note"`
}

type CountRes struct {
	Data       []Expense `json:"data"`
	Page       int       `json:"page"`
	TotalPages int       `json:"total_pages"`
	Total      int       `json:"total"`
}

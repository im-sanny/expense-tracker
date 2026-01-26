package model

import "time"

type Expense struct {
	ID     int       `json:"id" db:"id"`
	Date   time.Time `json:"date" db:"date"`
	Amount *int64     `json:"amount" db:"amount"`
	Note   *string    `json:"note" db:"note"`
}

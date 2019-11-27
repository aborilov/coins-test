package payment

import "time"

// Transaction model
type Transaction struct {
	ID     int64     `json:"id"`
	From   int64     `json:"from"`
	To     int64     `json:"to"`
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
}

// Balance model
type Balance struct {
	AccountID int64   `json:"account_id"`
	Balance   float64 `json:"balance"`
}

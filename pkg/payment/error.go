package payment

import "fmt"

// ErrInsufficientFunds raised when account doesn't have sufficient funds
type ErrInsufficientFunds struct {
	ID int64
}

func (e ErrInsufficientFunds) Error() string {
	return fmt.Sprintf("insufficient funds, account with ID %d", e.ID)
}

package account

import "fmt"

// ErrNotFound - raised when account not found
type ErrNotFound struct {
	ID int64
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("account with ID %d not found", e.ID)
}

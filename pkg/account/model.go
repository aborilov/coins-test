package account

// Account model
type Account struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// New create new account model
func New(firstName, lastName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
	}
}

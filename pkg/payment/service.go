package payment

import (
	"coins/pkg/account"
	"context"
	"time"
)

// Service interface
type Service interface {
	GetBalance(context.Context, *account.Account) (*Balance, error)
	ListTransactions(context.Context, *account.Account) ([]*Transaction, error)
	Transfer(ctx context.Context, from, to *account.Account, amount float64) (*Transaction, error)
	TopUp(ctx context.Context, a *account.Account, amount float64) (*Balance, error)
}

// Repository interface
type Repository interface {
	GetBalance(ctx context.Context, accountID int64) (*Balance, error)
	ListTransactions(ctx context.Context, accountID int64) ([]*Transaction, error)
	Transfer(context.Context, *Transaction) (*Transaction, error)
	TopUp(ctx context.Context, accountID int64, amount float64) (*Balance, error)
}

type service struct {
	repo Repository
}

// TopUp - add funds to account balance
func (s *service) TopUp(ctx context.Context, a *account.Account, amount float64) (*Balance, error) {
	return s.repo.TopUp(ctx, a.ID, amount)
}

// GetBalance - return account balance or ErrNotFound if such account doesn't exists
func (s *service) GetBalance(ctx context.Context, a *account.Account) (*Balance, error) {
	return s.repo.GetBalance(ctx, a.ID)
}

// ListTransactions - return list of all account transactions or ErrNotFound if such account doesn't exists
func (s *service) ListTransactions(ctx context.Context, a *account.Account) ([]*Transaction, error) {
	return s.repo.ListTransactions(ctx, a.ID)
}

// Transfer -  transfer funds from one account to an other, raise ErrInsufficientFunds when `from` account doesn't have enough funds
func (s *service) Transfer(ctx context.Context, from, to *account.Account, amount float64) (*Transaction, error) {
	t := &Transaction{
		From:   from.ID,
		To:     to.ID,
		Amount: amount,
		Date:   time.Now().UTC(),
	}
	return s.repo.Transfer(ctx, t)
}

// NewService - build new service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

package account

import "context"

// Repository interface
type Repository interface {
	List(context.Context) ([]*Account, error)
	Get(ctx context.Context, id int64) (*Account, error)
	Store(context.Context, *Account) (*Account, error)
}

// Service interface
type Service interface {
	List(context.Context) ([]*Account, error)
	Get(ctx context.Context, id int64) (*Account, error)
	Store(ctx context.Context, firstName, lastName string) (*Account, error)
}

type service struct {
	repo Repository
}

// List return list of accounts without paging
func (s *service) List(ctx context.Context) ([]*Account, error) {
	return s.repo.List(ctx)
}

// Get return Account with requested id or return ErrNotFound error
func (s *service) Get(ctx context.Context, id int64) (*Account, error) {
	return s.repo.Get(ctx, id)
}

// Store new account with providerd first and lastname
func (s *service) Store(ctx context.Context, firstName, lastName string) (*Account, error) {
	return s.repo.Store(ctx, New(firstName, lastName))
}

// NewService build new Service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

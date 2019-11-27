package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeAddAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addAccountRequest)
		account, err := s.Store(ctx, req.FirstName, req.LastName)
		return addAccountResponse{Account: account, Err: err}, nil
	}
}

type addAccountRequest struct {
	FirstName string
	LastName  string
}

type addAccountResponse struct {
	Account *Account `json:"account,omitempty"`
	Err     error    `json:"error,omitempty"`
}

func (r addAccountResponse) error() error { return r.Err }

func makeGetAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAccountRequest)
		account, err := s.Get(ctx, req.ID)
		return getAccountResponse{Account: account, Err: err}, err
	}
}

type getAccountRequest struct {
	ID int64
}

type getAccountResponse struct {
	Account *Account `json:"account,omitempty"`
	Err     error    `json:"err,omitempty"`
}

func (r getAccountResponse) error() error { return r.Err }

func makeListAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(listAccountsRequest)
		accounts, err := s.List(ctx)
		return listAccountsResponse{Accounts: accounts, Err: err}, err
	}
}

type listAccountsRequest struct{}

type listAccountsResponse struct {
	Accounts []*Account `json:"accounts,omitempty"`
	Err      error      `json:"err,omitempty"`
}

func (r listAccountsResponse) error() error { return r.Err }

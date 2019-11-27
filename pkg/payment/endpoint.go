package payment

import (
	"coins/pkg/account"
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeGetBalanceEndpoint(s Service, as account.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getBalanceRequest)

		a, err := as.Get(ctx, req.ID)
		if err != nil {
			return getBalanceResponse{Err: err}, err
		}
		balance, err := s.GetBalance(ctx, a)
		return getBalanceResponse{Balance: balance, Err: err}, err
	}
}

type getBalanceRequest struct {
	ID int64
}

type getBalanceResponse struct {
	Balance *Balance `json:"balance,omitempty"`
	Err     error    `json:"err,omitempty"`
}

func (r getBalanceResponse) error() error { return r.Err }

func makeListTransactionsEndpoint(s Service, as account.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listTransactionsRequest)
		a, err := as.Get(ctx, req.ID)
		if err != nil {
			return listTransactionsResponse{Err: err}, err
		}

		tt, err := s.ListTransactions(ctx, a)
		return listTransactionsResponse{Transactions: tt, Err: err}, err
	}
}

type listTransactionsRequest struct {
	ID int64
}

type listTransactionsResponse struct {
	Transactions []*Transaction `json:"transactions,omitempty"`
	Err          error          `json:"err,omitempty"`
}

func makeTransferEndpoint(s Service, as account.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transferRequest)
		from, err := as.Get(ctx, req.From)
		if err != nil {
			return transferResponse{Err: err}, err
		}
		to, err := as.Get(ctx, req.To)
		if err != nil {
			return transferResponse{Err: err}, err
		}

		t, err := s.Transfer(ctx, from, to, req.Amount)
		return transferResponse{Transaction: t, Err: err}, err
	}
}

type transferRequest struct {
	From   int64
	To     int64
	Amount float64
}

type transferResponse struct {
	Transaction *Transaction `json:"transaction,omitempty"`
	Err         error        `json:"err,omitempty"`
}

func makeTopUpEndpoint(s Service, as account.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(topUpRequest)
		a, err := as.Get(ctx, req.AccountID)
		if err != nil {
			return topUpResponse{Err: err}, err
		}

		b, err := s.TopUp(ctx, a, req.Amount)
		return topUpResponse{Balance: b, Err: err}, err
	}
}

type topUpRequest struct {
	AccountID int64
	Amount    float64
}

type topUpResponse struct {
	Balance *Balance `json:"balance,omitempty"`
	Err     error    `json:"err,omitempty"`
}

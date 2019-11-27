package payment

import (
	"coins/pkg/account"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type errBadRequest struct {
	Msg string
}

func (e errBadRequest) Error() string {
	return e.Msg
}

// MakeHandler build handlers for payment transport
func MakeHandler(ps Service, as account.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	getBalanceHandler := kithttp.NewServer(
		makeGetBalanceEndpoint(ps, as),
		decodeGetBalanceRequest,
		encodeResponse,
		opts...,
	)

	listTransactionsHandler := kithttp.NewServer(
		makeListTransactionsEndpoint(ps, as),
		decodeListTransactionsRequest,
		encodeResponse,
		opts...,
	)

	transferHandler := kithttp.NewServer(
		makeTransferEndpoint(ps, as),
		decodeTransferRequest,
		encodeResponse,
		opts...,
	)

	topUpHandler := kithttp.NewServer(
		makeTopUpEndpoint(ps, as),
		decodeTopUpRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/payment/v1/balance/{id}", getBalanceHandler).Methods("GET")
	r.Handle("/payment/v1/transactions/{id}", listTransactionsHandler).Methods("GET")
	r.Handle("/payment/v1/transfer", transferHandler).Methods("POST")
	r.Handle("/payment/v1/topup", topUpHandler).Methods("POST")

	return r
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err.(type) {
	case account.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case errBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	case ErrInsufficientFunds:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func decodeGetBalanceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, errBadRequest{Msg: fmt.Sprintf("id param required")}
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, errBadRequest{Msg: fmt.Sprintf("id param must be int")}
	}
	return getBalanceRequest{ID: id}, nil
}

func decodeListTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, errBadRequest{Msg: fmt.Sprintf("id param required")}
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, errBadRequest{Msg: fmt.Sprintf("id param must be int")}
	}
	return listTransactionsRequest{ID: id}, nil
}

func decodeTransferRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		From   int64   `json:"from"`
		To     int64   `json:"to"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return transferRequest{From: body.From, To: body.To, Amount: body.Amount}, nil
}

func decodeTopUpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		AccountID int64   `json:"account_id"`
		Amount    float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return topUpRequest{AccountID: body.AccountID, Amount: body.Amount}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

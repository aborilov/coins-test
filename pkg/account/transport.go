package account

import (
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

// MakeHandler build handlers for account transport
func MakeHandler(as Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	addAccountHandler := kithttp.NewServer(
		makeAddAccountEndpoint(as),
		decodeAddAccountRequest,
		encodeResponse,
		opts...,
	)

	getAccountHandler := kithttp.NewServer(
		makeGetAccountEndpoint(as),
		decodeGetAccountRequest,
		encodeResponse,
		opts...,
	)

	listAccountsHandler := kithttp.NewServer(
		makeListAccountsEndpoint(as),
		decodeListAccountsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/account/v1/", addAccountHandler).Methods("POST")
	r.Handle("/account/v1/{id}", getAccountHandler).Methods("GET")
	r.Handle("/account/v1/", listAccountsHandler).Methods("GET")

	return r
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err.(type) {
	case ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case errBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func decodeAddAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return addAccountRequest{
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}, nil
}

func decodeGetAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, errBadRequest{Msg: fmt.Sprintf("id param required")}
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, errBadRequest{Msg: fmt.Sprintf("id param must be int")}
	}
	return getAccountRequest{ID: id}, nil
}

func decodeListAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listAccountsRequest{}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

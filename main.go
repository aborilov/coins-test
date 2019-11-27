package main

import (
	"coins/pkg/account"
	"coins/pkg/payment"
	accountRepo "coins/repository/account/pg"
	paymentRepo "coins/repository/payment/pg"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/lib/pq"
)

const (
	defaultDbURI = "postgres://server:server@127.0.0.1:5432/server?sslmode=disable"
)

func getDB() *sql.DB {
	dbURI := os.Getenv("PG_URI")
	if dbURI == "" {
		dbURI = defaultDbURI
	}
	uri, err := pq.ParseURL(dbURI)
	if err != nil {
		panic(err)
	}
	pdb, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	return pdb
}
func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	as := account.NewService(accountRepo.NewRepository(getDB()))
	ps := payment.NewService(paymentRepo.NewRepository(getDB()))

	mux := http.NewServeMux()

	mux.Handle("/account/v1/", account.MakeHandler(as))
	mux.Handle("/payment/v1/", payment.MakeHandler(ps, as))

	http.Handle("/", mux)

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", ":80", "msg", "listening")
		errs <- http.ListenAndServe(":80", nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

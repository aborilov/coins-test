package pg

import (
	"coins/pkg/payment"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v8"
	_ "github.com/doug-martin/goqu/v8/dialect/postgres"
	"github.com/pkg/errors"
)

const (
	tableTransaction = "transaction"
	tableBalance     = "balance"
)

type recordBalance struct {
	AccountID int64   `db:"account_id"`
	Balance   float64 `db:"balance"`
}

func (b *recordBalance) toBalance() *payment.Balance {
	return &payment.Balance{
		AccountID: b.AccountID,
		Balance:   b.Balance,
	}
}

func fromBalance(b *payment.Balance) *recordBalance {
	return &recordBalance{
		AccountID: b.AccountID,
		Balance:   b.Balance,
	}
}

type recordTransaction struct {
	ID     int64     `db:"id" goqu:"skipinsert,skipupdate"`
	From   int64     `db:"from"`
	To     int64     `db:"to"`
	Amount float64   `db:"amount"`
	Date   time.Time `db:"date"`
}

func (t *recordTransaction) toTransaction() *payment.Transaction {
	return &payment.Transaction{
		ID:     t.ID,
		From:   t.From,
		To:     t.To,
		Amount: t.Amount,
		Date:   t.Date,
	}
}

func fromTransaction(t *payment.Transaction) *recordTransaction {
	return &recordTransaction{
		ID:     t.ID,
		From:   t.From,
		To:     t.To,
		Amount: t.Amount,
		Date:   t.Date,
	}
}

type repository struct {
	gq *goqu.Database
}

// NewRepository - build new repository
func NewRepository(db *sql.DB) payment.Repository {
	return &repository{gq: goqu.New("postgres", db)}
}

func (repo *repository) GetBalance(ctx context.Context, id int64) (*payment.Balance, error) {
	r := &recordBalance{}
	tx, err := repo.gq.Begin()
	if err != nil {
		return nil, err
	}
	r, err = getBalance(ctx, tx, id)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get balance")
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return r.toBalance(), nil
}

func (repo *repository) ListTransactions(ctx context.Context, id int64) ([]*payment.Transaction, error) {
	var rr []*recordTransaction
	if err := repo.gq.From(tableTransaction).Where(goqu.ExOr{"from": id, "to": id}).Order(goqu.I("date").Asc()).ScanStructsContext(ctx, &rr); err != nil {
		return nil, errors.Wrap(err, "unable to retrieve transaction records")
	}
	tt := make([]*payment.Transaction, 0, len(rr))
	for _, r := range rr {
		tt = append(tt, r.toTransaction())
	}

	return tt, nil
}

func lockBalance(tx *goqu.TxDatabase, accountID int64) error {
	res := tx.Select(goqu.L(fmt.Sprintf("pg_try_advisory_xact_lock(%d)", accountID))).Executor()

	var lock bool
	if _, err := res.ScanVal(&lock); err != nil {
		return errors.Wrapf(err, "unable to acquire lock for account with ID %d", accountID)
	}
	if !lock {
		return errors.Errorf("unable to acquire lock for account with ID %d", accountID)
	}
	return nil
}

func updateBalance(ctx context.Context, tx *goqu.TxDatabase, balance *recordBalance) error {
	_, err := tx.Insert(tableBalance).Rows(balance).OnConflict(goqu.DoUpdate("account_id", balance)).Executor().ExecContext(ctx)
	return err
}

func getBalance(ctx context.Context, tx *goqu.TxDatabase, id int64) (*recordBalance, error) {
	r := &recordBalance{}
	found, err := tx.From(tableBalance).Where(goqu.I("account_id").Eq(id)).ScanStructContext(ctx, r)
	if err != nil {
		return nil, err
	}
	r.AccountID = id
	if !found {
		return r, nil
	}
	return r, nil
}

func (repo *repository) Transfer(ctx context.Context, t *payment.Transaction) (*payment.Transaction, error) {

	r := fromTransaction(t)
	err := repo.gq.WithTx(func(tx *goqu.TxDatabase) error {

		if err := lockBalance(tx, t.From); err != nil {
			return err
		}

		if err := lockBalance(tx, t.To); err != nil {
			return err
		}

		fromBalance, err := getBalance(ctx, tx, t.From)
		if err != nil {
			return err
		}
		if fromBalance.Balance < t.Amount {
			return payment.ErrInsufficientFunds{ID: t.From}
		}

		fromBalance.Balance -= t.Amount
		if err := updateBalance(ctx, tx, fromBalance); err != nil {
			return err
		}

		toBalance, err := getBalance(ctx, tx, t.To)
		if err != nil {
			return err
		}
		toBalance.Balance += t.Amount
		if err := updateBalance(ctx, tx, toBalance); err != nil {
			return err
		}

		res := tx.From(tableTransaction).Insert().Returning(goqu.C("id")).Rows(r).Executor()
		var id int64
		if _, err := res.ScanVal(&id); err != nil {
			return errors.Wrap(err, "failed to retrieve last inserted ID")
		}
		t.ID = id
		return nil
	})
	return t, err
}

func (repo *repository) TopUp(ctx context.Context, accountID int64, amount float64) (*payment.Balance, error) {
	var b *recordBalance
	err := repo.gq.WithTx(func(tx *goqu.TxDatabase) error {
		if err := lockBalance(tx, accountID); err != nil {
			return err
		}
		var err error
		b, err = getBalance(ctx, tx, accountID)
		if err != nil {
			return err
		}
		b.Balance += amount
		if err := updateBalance(ctx, tx, b); err != nil {
			return err
		}
		return nil
	})
	return b.toBalance(), err
}

package pg

import (
	"coins/pkg/account"
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v8"
	_ "github.com/doug-martin/goqu/v8/dialect/postgres"
	"github.com/pkg/errors"
)

const table = "account"

type record struct {
	ID        int64  `db:"id" goqu:"skipinsert,skipupdate"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

func (t *record) toAccount() *account.Account {
	return &account.Account{
		ID:        t.ID,
		FirstName: t.FirstName,
		LastName:  t.LastName,
	}
}

func fromAccount(a *account.Account) *record {
	return &record{
		ID:        a.ID,
		FirstName: a.FirstName,
		LastName:  a.LastName,
	}
}

type repository struct {
	gq *goqu.Database
}

// NewRepository - build new repository
func NewRepository(db *sql.DB) account.Repository {
	return &repository{gq: goqu.New("postgres", db)}
}

func (repo *repository) List(ctx context.Context) ([]*account.Account, error) {
	var rr []*record
	if err := repo.gq.From(table).Order(goqu.I("id").Asc()).ScanStructsContext(ctx, &rr); err != nil {
		return nil, errors.Wrap(err, "unable to retrieve account records")
	}
	aa := make([]*account.Account, 0, len(rr))
	for _, r := range rr {
		aa = append(aa, r.toAccount())
	}

	return aa, nil
}
func (repo *repository) Get(ctx context.Context, id int64) (*account.Account, error) {
	r := &record{}
	found, err := repo.gq.From(table).Where(goqu.I("id").Eq(id)).ScanStructContext(ctx, r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get account")
	}
	if !found {
		return nil, account.ErrNotFound{ID: id}
	}
	return r.toAccount(), nil
}
func (repo *repository) Store(ctx context.Context, a *account.Account) (*account.Account, error) {
	r := fromAccount(a)
	res := repo.gq.From(table).Insert().Returning(goqu.C("id")).Rows(r).Executor()
	var id int64
	if _, err := res.ScanVal(&id); err != nil {
		return nil, errors.Wrap(err, "failed to retrieve last inserted ID")
	}
	a.ID = id
	return a, nil
}

package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/doug-martin/goqu/v9"
)

const (
	TableAccountPassword = "account_passwords"
	ColOfAccountID       = "of_account_id"
	ColHash              = "hash"
)

type AccountPassword struct {
	OfAccountID uint64 `sql:"of_account_id"`
	Hash        string `sql:"hash"`
}

type AccountPasswordDataAccessor interface {
	CreateAccountPassword(ctx context.Context, accountID uint64, passwordHash string) error
	UpdateAccountPassword(ctx context.Context, account AccountPassword) error
	GetAccountPasswordByAccountID(ctx context.Context, accountID uint64) (AccountPassword, error)
	WithDatabase(database Database) AccountPasswordDataAccessor
}

type accountPasswordAccessor struct {
	database Database
}

func NewAccountPasswordDataAccessor(database *goqu.Database) AccountPasswordDataAccessor {
	return &accountPasswordAccessor{database: database}
}

func (a accountPasswordAccessor) UpdateAccountPassword(ctx context.Context, account AccountPassword) error {
	_, err := a.database.Update(TableAccountPassword).Set(goqu.Record{
		ColHash: account.Hash,
	}).Where(goqu.Ex{ColOfAccountID: account.OfAccountID}).Executor().ExecContext(ctx)
	if err != nil {
		log.Print("error updating account password:", err)
		return err
	}
	return nil
}

func (a accountPasswordAccessor) CreateAccountPassword(ctx context.Context, accountID uint64, passwordHash string) error {
	_, err := a.database.Insert(TableAccountPassword).Rows(goqu.Record{
		ColOfAccountID: accountID,
		ColHash:        passwordHash,
	}).Executor().ExecContext(ctx)
	if err != nil {
		log.Print("error inserting account password:", err)
		return err
	}
	return nil
}

func (a accountPasswordAccessor) WithDatabase(database Database) AccountPasswordDataAccessor {
	return &accountPasswordAccessor{
		database: database,
	}
}

func (a accountPasswordAccessor) GetAccountPasswordByAccountID(ctx context.Context, accountID uint64) (AccountPassword, error) {
	var password AccountPassword
	found, err := a.database.From(TableAccountPassword).
		Select(ColOfAccountID, ColHash).
		Where(goqu.Ex{ColOfAccountID: accountID}).
		ScanStructContext(ctx, &password)

	if err != nil {
		log.Print("error getting account password:", err)
		return AccountPassword{}, err
	}

	if !found {
		return AccountPassword{}, sql.ErrNoRows
	}

	return password, nil
}

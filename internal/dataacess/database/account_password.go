package database

import (
	"context"
	"log"

	"github.com/doug-martin/goqu/v9"
)

type AccountPasswordDataAccessor interface {
	CreateAccountPassword(ctx context.Context, accountID uint64, passwordHash string) error
	UpdateAccountPassword(ctx context.Context, accountID uint64, passwordHash string) error
    WithDatabase(database Database) AccountPasswordDataAccessor
}

type accountPasswordAccessor struct {
	database Database
}

func NewAccountPasswordDataAccessor(database *goqu.Database) AccountPasswordDataAccessor {
	return &accountPasswordAccessor{database: database}
}

func (a *accountPasswordAccessor) UpdateAccountPassword(ctx context.Context, accountID uint64, passwordHash string) error {
	panic("unimplemented")
}

func (a accountPasswordAccessor) CreateAccountPassword(ctx context.Context, accountID uint64, passwordHash string) error {
	_, err := a.database.Insert("account_passwords").Rows(goqu.Record{
		"of_account_id": accountID,
		"hash":          passwordHash,
	}).Executor().Exec()
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

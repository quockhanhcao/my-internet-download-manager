package database

import (
	"context"
	"log"

	"github.com/doug-martin/goqu/v9"
)

type Account struct {
	AccountID   uint64 `sql:"account_id"`
	AccountName string `sql:"account_name"`
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, accountName string, password string) (uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByAccountName(ctx context.Context, accountName string) (Account, error)
	WithDatabase(database Database) AccountDataAccessor
}
type accountDataAccessor struct {
	database Database
}

func NewAccountDataAccessor(database *goqu.Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
	}
}

func (a accountDataAccessor) CreateAccount(ctx context.Context, accountName string, password string) (uint64, error) {
	result, err := a.database.Insert("accounts").Rows(goqu.Record{
		"account_name": accountName,
	}).Executor().Exec()
	if err != nil {
		log.Println("error inserting account:", err)
		return 0, err
	}
	accountID, err := result.LastInsertId()
	if err != nil {
		log.Println("error getting last insert ID:", err)
		return 0, err
	}
	return uint64(accountID), nil
}

func (a accountDataAccessor) GetAccountByAccountName(ctx context.Context, accountName string) (Account, error) {
	return Account{}, nil // Placeholder for actual implementation
}

func (a accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	panic("unimplemented")
}

func (a accountDataAccessor) WithDatabase(database Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
	}
}

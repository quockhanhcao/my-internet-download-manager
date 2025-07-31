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
	CreateAccount(ctx context.Context, account Account) (uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByAccountName(ctx context.Context, accountName string) (Account, error)
	WithDatabase(database Database) AccountDataAccessor
}

type accountAccessor struct {
	database Database
}

// CreateAccount implements AccountAccessor.
func (a *accountAccessor) CreateAccount(ctx context.Context, account Account) (uint64, error) {
	result, err := a.database.Insert("accounts").Rows(goqu.Record{
		"account_name": account.AccountName,
	}).Executor().ExecContext(ctx)
	if err != nil {
		log.Printf("error creating account: %v", err)
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("error getting last insert ID: %v", err)
		return 0, err
	}
	return uint64(lastInsertID), nil
}

// GetAccountByID implements AccountAccessor.
func (a *accountAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	panic("unimplemented")
}

// GetAccountByAccountName implements AccountAccessor.
func (a *accountAccessor) GetAccountByAccountName(ctx context.Context, accountName string) (Account, error) {
	panic("unimplemented")
}

func NewAccountDataAccessor(database *goqu.Database) AccountDataAccessor {
	return &accountAccessor{database}
}

func (a *accountAccessor) WithDatabase(database Database) AccountDataAccessor {
	return &accountAccessor{database: database}
}

package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

type AccountPassword struct {
	OfAccountID uint64 `sql:"of_account_id"`
	Hash        string `sql:"hash"`
}

type AccountPasswordDataAccessor interface {
	CreateAccountPassword(ctx context.Context, accountPassword AccountPassword) error
	GetAccountPasswordByAccountID(ctx context.Context, ofAccountID uint64) (AccountPassword, error)
	UpdateAccountPassword(ctx context.Context, ofAccountID uint64, hash string) error
    WithDatabase(database Database) AccountPasswordDataAccessor
}

type accountPasswordAccessor struct {
	database Database
}

// UpdateAccountPassword implements AccountPasswordDataAccessor.
func (a *accountPasswordAccessor) UpdateAccountPassword(ctx context.Context, ofAccountID uint64, hash string) error {
	panic("unimplemented")
}

// CreateAccountPassword implements AccountPasswordDataAccessor.
func (a *accountPasswordAccessor) CreateAccountPassword(ctx context.Context, accountPassword AccountPassword) error {
	panic("unimplemented")
}

// GetAccountPasswordByAccountID implements AccountPasswordDataAccessor.
func (a *accountPasswordAccessor) GetAccountPasswordByAccountID(ctx context.Context, ofAccountID uint64) (AccountPassword, error) {
	panic("unimplemented")
}

func (a *accountPasswordAccessor) WithDatabase(database Database) AccountPasswordDataAccessor {
	return &accountPasswordAccessor{database: database}
}

func NewAccountPasswordDataAccessor(database *goqu.Database) AccountPasswordDataAccessor {
	return &accountPasswordAccessor{database: database}
}

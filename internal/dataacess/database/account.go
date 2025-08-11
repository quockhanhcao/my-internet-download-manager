package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

const (
	TableAccount = "accounts"
)

type Account struct {
	ID          uint64 `sql:"id"`
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
	logger   *zap.Logger
}

func NewAccountDataAccessor(database *goqu.Database, logger *zap.Logger) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (a accountDataAccessor) CreateAccount(ctx context.Context, accountName string, password string) (uint64, error) {
	a.logger.With(zap.String("accountName", accountName)).Info("creating account in database")

	result, err := a.database.Insert(TableAccount).Rows(goqu.Record{
		"account_name": accountName,
	}).Executor().ExecContext(ctx)
	if err != nil {
		a.logger.With(zap.Error(err), zap.String("accountName", accountName)).Error("failed to insert account")
		return 0, err
	}

	accountID, err := result.LastInsertId()
	if err != nil {
		a.logger.With(zap.Error(err), zap.String("accountName", accountName)).Error("failed to get last insert ID")
		return 0, err
	}

	a.logger.With(zap.Uint64("accountID", uint64(accountID)), zap.String("accountName", accountName)).Info("account created successfully in database")
	return uint64(accountID), nil
}

func (a accountDataAccessor) GetAccountByAccountName(ctx context.Context, accountName string) (Account, error) {
	a.logger.With(zap.String("accountName", accountName)).Info("getting account by account name")

	var account Account
	found, err := a.database.From(TableAccount).
		Select("id", "account_name").
		Where(goqu.Ex{"account_name": accountName}).
		ScanStructContext(ctx, &account)

	if err != nil {
		a.logger.With(zap.Error(err), zap.String("accountName", accountName)).Error("failed to get account by account name")
		return Account{}, err
	}

	if !found {
		a.logger.With(zap.String("accountName", accountName)).Warn("account not found by account name")
		return Account{}, nil
	}

	a.logger.With(zap.Uint64("accountID", account.ID), zap.String("accountName", accountName)).Info("account found by account name")
	return account, nil
}

func (a accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	a.logger.With(zap.Uint64("accountID", id)).Info("getting account by ID")

	var account Account
	found, err := a.database.From(TableAccount).
		Select("id", "account_name").
		Where(goqu.Ex{"id": id}).
		ScanStructContext(ctx, &account)

	if err != nil {
		a.logger.With(zap.Error(err), zap.Uint64("accountID", id)).Error("failed to get account by ID")
		return Account{}, err
	}

	if !found {
		a.logger.With(zap.Uint64("accountID", id)).Warn("account not found by ID")
		return Account{}, nil // Return empty struct and nil error when not found
	}

	a.logger.With(zap.Uint64("accountID", account.ID), zap.String("accountName", account.AccountName)).Info("account found by ID")
	return account, nil
}

func (a accountDataAccessor) WithDatabase(database Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
		logger:   a.logger,
	}
}

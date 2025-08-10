package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
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

type accountPasswordDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewAccountPasswordDataAccessor(database *goqu.Database, logger *zap.Logger) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{database: database, logger: logger}
}

func (a accountPasswordDataAccessor) UpdateAccountPassword(ctx context.Context, account AccountPassword) error {
	a.logger.With(zap.Uint64("accountID", account.OfAccountID)).Info("updating account password")

	_, err := a.database.Update(TableAccountPassword).Set(goqu.Record{
		ColHash: account.Hash,
	}).Where(goqu.Ex{ColOfAccountID: account.OfAccountID}).Executor().ExecContext(ctx)
	if err != nil {
		a.logger.With(zap.Error(err), zap.Uint64("accountID", account.OfAccountID)).Error("failed to update account password")
		return err
	}

	a.logger.With(zap.Uint64("accountID", account.OfAccountID)).Info("account password updated successfully")
	return nil
}

func (a accountPasswordDataAccessor) CreateAccountPassword(ctx context.Context, accountID uint64, passwordHash string) error {
	a.logger.With(zap.Uint64("accountID", accountID)).Info("creating account password")

	_, err := a.database.Insert(TableAccountPassword).Rows(goqu.Record{
		ColOfAccountID: accountID,
		ColHash:        passwordHash,
	}).Executor().ExecContext(ctx)
	if err != nil {
		a.logger.With(zap.Error(err), zap.Uint64("accountID", accountID)).Error("failed to insert account password")
		return err
	}

	a.logger.With(zap.Uint64("accountID", accountID)).Info("account password created successfully")
	return nil
}

func (a accountPasswordDataAccessor) WithDatabase(database Database) AccountPasswordDataAccessor {
	return &accountPasswordDataAccessor{
		database: database,
	}
}

func (a accountPasswordDataAccessor) GetAccountPasswordByAccountID(ctx context.Context, accountID uint64) (AccountPassword, error) {
	a.logger.With(zap.Uint64("accountID", accountID)).Info("getting account password by account ID")

	var password AccountPassword
	found, err := a.database.From(TableAccountPassword).
		Select(ColOfAccountID, ColHash).
		Where(goqu.Ex{ColOfAccountID: accountID}).
		ScanStructContext(ctx, &password)

	if err != nil {
		a.logger.With(zap.Error(err), zap.Uint64("accountID", accountID)).Error("failed to get account password")
		return AccountPassword{}, err
	}

	if !found {
		a.logger.With(zap.Uint64("accountID", accountID)).Warn("account password not found")
		return AccountPassword{}, sql.ErrNoRows
	}

	a.logger.With(zap.Uint64("accountID", accountID)).Info("account password retrieved successfully")
	return password, nil
}

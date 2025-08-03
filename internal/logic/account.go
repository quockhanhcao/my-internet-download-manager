package logic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/quockhanhcao/my-internet-download-manager/internal/dataacess/database"
)

type Account struct {
	AccountID   uint64
	AccountName string
}

type CreateAccountParams struct {
	AccountName string
	Password    string
}

type CreateSessionParams struct {
	AccountName string
	Password    string
}

type CreateSessionResponse struct {
	Account Account
	Token   string
}

type AccountHandler interface {
	CreateAccount(ctx context.Context, params CreateAccountParams) (Account, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionResponse, error)
}

type accountHandler struct {
	accountDataAccessor         database.AccountDataAccessor
	accountPasswordDataAccessor database.AccountPasswordDataAccessor
	hash                        Hash
	goquDatabase                *goqu.Database
}

func NewAccountHandler(
	accountDataAccessor database.AccountDataAccessor,
	accountPasswordDataAccessor database.AccountPasswordDataAccessor,
	hash Hash,
	goquDatabase *goqu.Database,
) AccountHandler {
	return &accountHandler{
		accountDataAccessor:         accountDataAccessor,
		accountPasswordDataAccessor: accountPasswordDataAccessor,
		hash:                        hash,
		goquDatabase:                goquDatabase,
	}
}

func (a accountHandler) isAccountExisted(ctx context.Context, accountName string) (bool, error) {
	_, err := a.accountDataAccessor.GetAccountByAccountName(ctx, accountName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a accountHandler) CreateAccount(ctx context.Context, params CreateAccountParams) (Account, error) {
	var accountID uint64
	txErr := a.goquDatabase.WithTx(func(tx *goqu.TxDatabase) error {
		accountNameTaken, err := a.isAccountExisted(ctx, params.AccountName)
		if err != nil {
			return err
		}
		if accountNameTaken {
			return errors.New("account name already taken")
		}

		accountID, err = a.accountDataAccessor.WithDatabase(tx).CreateAccount(ctx, params.AccountName, params.Password)
		if err != nil {
			return err
		}

		hashedPassword, err := a.hash.HashPassword(params.Password)
		if err != nil {
			return err
		}
		a.accountPasswordDataAccessor.WithDatabase(tx).CreateAccountPassword(ctx, accountID, hashedPassword)
		return nil
	})
	if txErr != nil {
		return Account{}, txErr
	}
	return Account{
		AccountID:   accountID,
		AccountName: params.AccountName,
	}, nil
}

func (a accountHandler) CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionResponse, error) {
	panic("unimplemented")
}

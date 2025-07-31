package logic

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/quockhanhcao/my-internet-download-manager/internal/dataacess/database"
)

type CreateAccountParams struct {
	AccountName string
	Password    string
}

type Account struct {
	ID          uint64
	AccountName string
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
	accountDataAccessor     database.AccountDataAccessor
	accountPasswordAccessor database.AccountPasswordDataAccessor
	hash                    Hash
	goquDatabase            *goqu.Database
}

func NewAccountHandler(
	accountDataAccessor database.AccountDataAccessor,
	accountPasswordAccessor database.AccountPasswordDataAccessor,
	hash Hash,
	goquDatabase *goqu.Database,
) AccountHandler {
	return &accountHandler{
		goquDatabase:            goquDatabase,
		accountDataAccessor:     accountDataAccessor,
		accountPasswordAccessor: accountPasswordAccessor,
		hash:                    hash,
	}
}

func (a accountHandler) isAccountNameTaken(ctx context.Context, accountName string) (bool, error) {
	_, err := a.accountDataAccessor.GetAccountByAccountName(ctx, accountName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		log.Printf("error checking account name: %v", err)
		return false, err
	}
	return true, nil
}

func (a accountHandler) CreateAccount(ctx context.Context, params CreateAccountParams) (Account, error) {
	var accountID uint64
	txErr := a.goquDatabase.WithTx(func(td *goqu.TxDatabase) error {
		accountNameTaken, err := a.isAccountNameTaken(ctx, params.AccountName)
		if err != nil {
			return err
		}
		if accountNameTaken {
			log.Printf("account name %s is already taken", params.AccountName)
			return errors.New("account name is already taken")
		}
		// insert into account table
		accountID, err = a.accountDataAccessor.WithDatabase(td).CreateAccount(ctx, database.Account{AccountName: params.AccountName})
		if err != nil {
			log.Printf("error creating account: %v", err)
			return err
		}
		hashPassword, err := a.hash.Hash(ctx, params.Password)
		if err != nil {
			return err
		}
		err = a.accountPasswordAccessor.WithDatabase(td).CreateAccountPassword(ctx, database.AccountPassword{
			OfAccountID: accountID,
			Hash:        string(hashPassword),
		})
		if err != nil {
			log.Printf("error creating account password: %v", err)
			return err
		}
		return nil
	})
	if txErr != nil {
		log.Printf("transaction error: %v", txErr)
		return Account{}, txErr
	}
	return Account{
		ID:          accountID,
		AccountName: params.AccountName,
	}, nil
}

func (a accountHandler) CreateSession(ctx context.Context, params CreateSessionParams) (CreateSessionResponse, error) {
	// Implementation for creating a session
	return CreateSessionResponse{}, nil
}

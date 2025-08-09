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

type AccountHandler interface {
	CreateAccount(ctx context.Context, params CreateAccountParams) (Account, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (token string, err error)
}

type accountHandler struct {
	accountDataAccessor         database.AccountDataAccessor
	accountPasswordDataAccessor database.AccountPasswordDataAccessor
	tokenPublicKeyDataAccessor  database.TokenPublicKeyDataAccessor
	hashHandler                 HashHandler
	tokenHandler                TokenHandler
	goquDatabase                *goqu.Database
}

func NewAccountHandler(
	accountDataAccessor database.AccountDataAccessor,
	accountPasswordDataAccessor database.AccountPasswordDataAccessor,
	tokenPublicKeyDataAccessor database.TokenPublicKeyDataAccessor,
	hashHandler HashHandler,
	tokenHandler TokenHandler,
	goquDatabase *goqu.Database,
) AccountHandler {
	return &accountHandler{
		accountDataAccessor:         accountDataAccessor,
		accountPasswordDataAccessor: accountPasswordDataAccessor,
		tokenPublicKeyDataAccessor:  tokenPublicKeyDataAccessor,
		hashHandler:                 hashHandler,
		tokenHandler:                tokenHandler,
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

		hashedPassword, err := a.hashHandler.Hash(ctx, params.Password)
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

func (a accountHandler) CreateSession(ctx context.Context, params CreateSessionParams) (token string, err error) {
	existingAccount, err := a.accountDataAccessor.GetAccountByAccountName(ctx, params.AccountName)
	if err != nil {
		return "", err
	}
	existingPassword, err := a.accountPasswordDataAccessor.GetAccountPasswordByAccountID(ctx, existingAccount.ID)
	if err != nil {
		return "", err
	}
	isHashEqual, err := a.hashHandler.IsHashEqual(ctx, existingPassword.Hash, params.Password)
	if err != nil {
		return "", err
	}

	if !isHashEqual {
		return "", errors.New("incorrect password")
	}

	// // generate a token
	// token, expiresIn, err := a.tokenHandler.GetToken(ctx, existingAccount.ID)
	// if err != nil {
	// 	return "", err
	// }
	return "", nil
}

package logic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/quockhanhcao/my-internet-download-manager/internal/dataacess/cache"
	"github.com/quockhanhcao/my-internet-download-manager/internal/dataacess/database"
	"go.uber.org/zap"
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
	logger                      *zap.Logger
	accountNameCache            cache.AccountNameCache
}

func NewAccountHandler(
	accountDataAccessor database.AccountDataAccessor,
	accountPasswordDataAccessor database.AccountPasswordDataAccessor,
	tokenPublicKeyDataAccessor database.TokenPublicKeyDataAccessor,
	hashHandler HashHandler,
	tokenHandler TokenHandler,
	goquDatabase *goqu.Database,
	logger *zap.Logger,
	accountNameCache cache.AccountNameCache,
) AccountHandler {
	return &accountHandler{
		accountDataAccessor:         accountDataAccessor,
		accountPasswordDataAccessor: accountPasswordDataAccessor,
		tokenPublicKeyDataAccessor:  tokenPublicKeyDataAccessor,
		hashHandler:                 hashHandler,
		tokenHandler:                tokenHandler,
		goquDatabase:                goquDatabase,
		logger:                      logger,
		accountNameCache:            accountNameCache,
	}
}

func (a accountHandler) isAccountExisted(ctx context.Context, accountName string) (bool, error) {
	cachedAccountName, err := a.accountNameCache.IsAccountNameTaken(ctx, accountName)
	if cachedAccountName {
		return true, nil
	}
	if err != nil {
		a.logger.With(zap.Error(err)).Warn("failed to check account name in cache")
	}
	_, err = a.accountDataAccessor.GetAccountByAccountName(ctx, accountName)
	if err != nil {
		return false, err
	}
	err = a.accountNameCache.SetAccountNameTaken(ctx, accountName)
	if err != nil {
		a.logger.With(zap.Error(err)).Warn("failed to set account name in cache")
	}
	return true, nil
}

func (a accountHandler) CreateAccount(ctx context.Context, params CreateAccountParams) (Account, error) {
	a.logger.With(zap.String("accountName", params.AccountName)).Info("starting account creation")

	var accountID uint64
	txErr := a.goquDatabase.WithTx(func(tx *goqu.TxDatabase) error {
		accountNameTaken, err := a.isAccountExisted(ctx, params.AccountName)
		if err != nil {
			a.logger.With(zap.Error(err)).Error("failed to check if account name is taken")
			return err
		}
		if accountNameTaken {
			a.logger.With(zap.String("accountName", params.AccountName)).Error("account name already taken")
			return errors.New("account name already taken")
		}

		accountID, err = a.accountDataAccessor.WithDatabase(tx).CreateAccount(ctx, params.AccountName, params.Password)
		if err != nil {
			a.logger.With(zap.Error(err), zap.String("accountName", params.AccountName)).Error("failed to create account in database")
			return err
		}
		a.logger.With(zap.Uint64("accountID", accountID), zap.String("accountName", params.AccountName)).Info("account created successfully")

		hashedPassword, err := a.hashHandler.Hash(ctx, params.Password)
		if err != nil {
			a.logger.With(zap.Error(err), zap.Uint64("accountID", accountID)).Error("failed to hash password")
			return err
		}

		err = a.accountPasswordDataAccessor.WithDatabase(tx).CreateAccountPassword(ctx, accountID, hashedPassword)
		if err != nil {
			a.logger.With(zap.Error(err), zap.Uint64("accountID", accountID)).Error("failed to create account password")
			return err
		}
		a.logger.With(zap.Uint64("accountID", accountID)).Info("account password created successfully")

		return nil
	})
	if txErr != nil {
		a.logger.With(zap.Error(txErr), zap.String("accountName", params.AccountName)).Error("account creation transaction failed")
		return Account{}, txErr
	}

	a.logger.With(zap.Uint64("accountID", accountID), zap.String("accountName", params.AccountName)).Info("account creation completed successfully")
	return Account{
		AccountID:   accountID,
		AccountName: params.AccountName,
	}, nil
}

func (a accountHandler) CreateSession(ctx context.Context, params CreateSessionParams) (token string, err error) {
	a.logger.With(zap.String("accountName", params.AccountName)).Info("starting session creation")

	existingAccount, err := a.accountDataAccessor.GetAccountByAccountName(ctx, params.AccountName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.logger.With(zap.String("accountName", params.AccountName)).Warn("account not found during session creation")
			return "", errors.New("account not found")
		}
		a.logger.With(zap.Error(err), zap.String("accountName", params.AccountName)).Error("failed to get account by name")
		return "", err
	}
	a.logger.With(zap.Uint64("accountID", existingAccount.ID), zap.String("accountName", params.AccountName)).Info("account found")

	existingPassword, err := a.accountPasswordDataAccessor.GetAccountPasswordByAccountID(ctx, existingAccount.ID)
	if err != nil {
		a.logger.With(zap.Error(err), zap.Uint64("accountID", existingAccount.ID)).Error("failed to get account password")
		return "", err
	}

	isHashEqual, err := a.hashHandler.IsHashEqual(ctx, existingPassword.Hash, params.Password)
	if err != nil {
		a.logger.With(zap.Error(err), zap.Uint64("accountID", existingAccount.ID)).Error("failed to verify password hash")
		return "", err
	}

	if !isHashEqual {
		a.logger.With(zap.Uint64("accountID", existingAccount.ID), zap.String("accountName", params.AccountName)).Warn("incorrect password provided")
		return "", errors.New("incorrect password")
	}
	a.logger.With(zap.Uint64("accountID", existingAccount.ID)).Info("password verified successfully")

	// generate a token
	token, _, err = a.tokenHandler.GetToken(ctx, existingAccount.ID)
	if err != nil {
		a.logger.With(zap.Error(err), zap.Uint64("accountID", existingAccount.ID)).Error("failed to generate token")
		return "", err
	}
	a.logger.With(zap.Uint64("accountID", existingAccount.ID), zap.String("accountName", params.AccountName)).Info("session created successfully")
	return token, nil
}

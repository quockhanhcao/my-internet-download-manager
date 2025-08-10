package database

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

const (
	TableTokenPublicKey   = "token_public_keys"
	ColOfTokenPublicKeyID = "id"
	ColPublicKey          = "public_key"
)

type TokenPublicKey struct {
	ID        uint64 `sql:"id"`
	PublicKey []byte `sql:"public_key"`
}

type TokenPublicKeyDataAccessor interface {
	CreateTokenPublicKey(ctx context.Context, publicKey []byte) (uint64, error)
	GetTokenPublicKeyByID(ctx context.Context, id uint64) (TokenPublicKey, error)
}

type tokenPublicKeyDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewTokenPublicKeyDataAccessor(database *goqu.Database, logger *zap.Logger) TokenPublicKeyDataAccessor {
	return &tokenPublicKeyDataAccessor{
		database: database,
		logger:   logger,
	}
}

// CreateTokenPublicKey implements TokenPublicKeyDataAccessor.
func (t tokenPublicKeyDataAccessor) CreateTokenPublicKey(ctx context.Context, publicKey []byte) (uint64, error) {
	t.logger.With(zap.Int("publicKeyLength", len(publicKey))).Info("creating token public key")

	result, err := t.database.Insert(TableTokenPublicKey).Rows(goqu.Record{
		ColPublicKey: publicKey,
	}).Executor().ExecContext(ctx)
	if err != nil {
		t.logger.With(zap.Error(err), zap.Int("publicKeyLength", len(publicKey))).Error("failed to insert token public key")
		return 0, err
	}
	keyID, err := result.LastInsertId()
	if err != nil {
		t.logger.With(zap.Error(err)).Error("failed to get last insert ID for token public key")
		return 0, err
	}

	t.logger.With(zap.Uint64("keyID", uint64(keyID)), zap.Int("publicKeyLength", len(publicKey))).Info("token public key created successfully")
	return uint64(keyID), nil
}

// GetTokenPublicKeyByID implements TokenPublicKeyDataAccessor.
func (t tokenPublicKeyDataAccessor) GetTokenPublicKeyByID(ctx context.Context, keyID uint64) (TokenPublicKey, error) {
	t.logger.With(zap.Uint64("keyID", keyID)).Info("getting token public key by ID")

	var publicKey TokenPublicKey
	found, err := t.database.From(TableTokenPublicKey).
		Select(ColOfTokenPublicKeyID, ColPublicKey).
		Where(goqu.Ex{ColOfTokenPublicKeyID: keyID}).
		ScanStructContext(ctx, &publicKey)

	if err != nil {
		t.logger.With(zap.Error(err), zap.Uint64("keyID", keyID)).Error("failed to get token public key")
		return TokenPublicKey{}, err
	}

	if !found {
		t.logger.With(zap.Uint64("keyID", keyID)).Warn("token public key not found")
		return TokenPublicKey{}, sql.ErrNoRows
	}

	t.logger.With(zap.Uint64("keyID", keyID), zap.Int("publicKeyLength", len(publicKey.PublicKey))).Info("token public key retrieved successfully")
	return publicKey, nil
}

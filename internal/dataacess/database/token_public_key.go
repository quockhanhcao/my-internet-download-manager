package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/doug-martin/goqu/v9"
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
}

func NewTokenPublicKeyDataAccessor(database *goqu.Database) TokenPublicKeyDataAccessor {
	return &tokenPublicKeyDataAccessor{
		database: database,
	}
}

// CreateTokenPublicKey implements TokenPublicKeyDataAccessor.
func (t tokenPublicKeyDataAccessor) CreateTokenPublicKey(ctx context.Context, publicKey []byte) (uint64, error) {
	result, err := t.database.Insert(TableTokenPublicKey).Rows(goqu.Record{
		ColPublicKey: publicKey,
	}).Executor().ExecContext(ctx)
	if err != nil {
		log.Println("error inserting account:", err)
		return 0, err
	}
	keyID, err := result.LastInsertId()
	if err != nil {
		log.Println("error getting last insert ID:", err)
		return 0, err
	}
	return uint64(keyID), nil
}

// GetTokenPublicKeyByID implements TokenPublicKeyDataAccessor.
func (t tokenPublicKeyDataAccessor) GetTokenPublicKeyByID(ctx context.Context, keyID uint64) (TokenPublicKey, error) {
	var publicKey TokenPublicKey
	found, err := t.database.From(TableTokenPublicKey).
		Select(ColOfTokenPublicKeyID, ColPublicKey).
		Where(goqu.Ex{ColOfTokenPublicKeyID: keyID}).
		ScanStructContext(ctx, &publicKey)

	if err != nil {
		log.Print("error getting token public key:", err)
		return TokenPublicKey{}, err
	}

	if !found {
		return TokenPublicKey{}, sql.ErrNoRows
	}

	return publicKey, nil
}

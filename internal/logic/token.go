package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"time"
)

const (
	rsaKeyPairBitSize = 2048
)

type TokenHandler interface {
	GetToken(ctx context.Context, accountID uint64) (string, time.Time, error)
	GetAccountIDAndExpireTime(ctx context.Context, token string) (uint64, time.Time, error)
}

type tokenHandler struct {
	privateKey *rsa.PrivateKey
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKeyPair, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	return privateKeyPair, nil
}

// GetAccountIDAndExpireTime implements TokenHandler.
func (t tokenHandler) GetAccountIDAndExpireTime(ctx context.Context, token string) (uint64, time.Time, error) {
	panic("unimplemented")
}

// GetToken implements TokenHandler.
func (t tokenHandler) GetToken(ctx context.Context, accountID uint64) (string, time.Time, error) {
	panic("unimplemented")
}

func NewTokenHandler() TokenHandler {
	privateKeyPair, err := generatePrivateKey(rsaKeyPairBitSize)
	if err != nil {
		return nil
	}
	return &tokenHandler{
		privateKey: privateKeyPair,
	}
}

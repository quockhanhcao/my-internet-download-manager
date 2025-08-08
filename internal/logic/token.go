package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/golang-jwt/jwt/v5"
	"github.com/quockhanhcao/my-internet-download-manager/internal/configs"
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
	publicKey  []byte
	authConfig configs.AuthConfig
}

func NewTokenHandler() TokenHandler {
	privateKeyPair, err := generatePrivateKey(rsaKeyPairBitSize)
	if err != nil {
		return nil
	}
	pemPublicKey, err := encodePublicKeyToPEM(privateKeyPair)
	if err != nil {
		return nil
	}

	return &tokenHandler{
		privateKey: privateKeyPair,
		publicKey:  pemPublicKey,
	}
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKeyPair, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	err = privateKeyPair.Validate()
	if err != nil {
		return nil, err
	}
	return privateKeyPair, nil
}

func encodePublicKeyToPEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	publicKey := &privateKey.PublicKey
	pubDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	pubBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubDER,
	}
	return pem.EncodeToMemory(&pubBlock), nil
}

// GetAccountIDAndExpireTime implements TokenHandler.
func (t tokenHandler) GetAccountIDAndExpireTime(ctx context.Context, token string) (uint64, time.Time, error) {
	panic("unimplemented")
}

// GetToken implements TokenHandler.
func (t tokenHandler) GetToken(ctx context.Context, accountID uint64) (string, time.Time, error) {
	expireTime := time.Now().Add(t.authConfig.ExpireTime)
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"kid": t.publicKey,
		"sub": accountID,
		"exp": expireTime,
	})

	signedToken, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", time.Time{}, err

	}
	return signedToken, expireTime, nil
}

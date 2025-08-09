package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/quockhanhcao/my-internet-download-manager/internal/configs"
	"github.com/quockhanhcao/my-internet-download-manager/internal/dataacess/database"
)

type TokenHandler interface {
	GetToken(ctx context.Context, accountID uint64) (string, time.Time, error)
	GetAccountIDAndExpireTime(ctx context.Context, token string) (uint64, time.Time, error)
    WithDatabase(database database.Database) TokenHandler
}

type tokenHandler struct {
	privateKey                 *rsa.PrivateKey
	configs                    configs.AuthConfig
	publicKeyID                uint64
	expiresIn                  time.Duration
	tokenPublicKeyDataAccessor database.TokenPublicKeyDataAccessor
    accountDataAccessor        database.AccountDataAccessor
}

func NewTokenHandler(
	configs configs.AuthConfig,
	tokenPublicKeyDataAccessor database.TokenPublicKeyDataAccessor,
	accountDataAccessor database.AccountDataAccessor,
) (TokenHandler, error) {
	expiresIn, err := configs.TokenConfig.GetExpiresInDuration()
	if err != nil {
		return nil, err
	}
	privateKeyPair, err := generateRSAKeyPair(configs.TokenConfig.KeyBitSize)
	if err != nil {
		return nil, err
	}
	pemPublicKey, err := encodePublicKeyToPEM(privateKeyPair)
	if err != nil {
		return nil, err
	}

	// save public key to database
	tokenPublicKeyID, err := tokenPublicKeyDataAccessor.CreateTokenPublicKey(context.Background(), pemPublicKey)
	if err != nil {
		return nil, err
	}

	return &tokenHandler{
		privateKey:                 privateKeyPair,
		publicKeyID:                tokenPublicKeyID,
		configs:                    configs,
		expiresIn:                  expiresIn,
		tokenPublicKeyDataAccessor: tokenPublicKeyDataAccessor,
        accountDataAccessor:        accountDataAccessor,
	}, nil
}

func generateRSAKeyPair(bitSize int) (*rsa.PrivateKey, error) {
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
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubDER,
	}
	return pem.EncodeToMemory(pubBlock), nil
}

func (t tokenHandler) decodePublicKeyFromPEM(ctx context.Context, keyID uint64) (*rsa.PublicKey, error) {
	pemPublicKey, err := t.tokenPublicKeyDataAccessor.GetTokenPublicKeyByID(ctx, keyID)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(pemPublicKey.PublicKey)
}

// GetAccountIDAndExpireTime implements TokenHandler.
func (t tokenHandler) GetAccountIDAndExpireTime(ctx context.Context, token string) (uint64, time.Time, error) {
    // parse the token
    parsedToken, err := jwt.Parse(token, func (parsedToken *jwt.Token) (interface{}, error) {
        // verify signing method is RSA
        if _, ok := parsedToken.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, jwt.ErrSignatureInvalid
        }

        // get key ID from token claims
        claims, ok := parsedToken.Claims.(jwt.MapClaims)
        if !ok {
            return nil, errors.New("can't get token claims")
        }

        tokenPublicKeyID, ok := claims["kid"].(float64)
        if !ok {
            return nil, errors.New("can't get token public key ID")
        }

        return t.decodePublicKeyFromPEM(ctx, uint64(tokenPublicKeyID))
    })

    if err != nil {
        return 0, time.Time{}, err
    }

    if !parsedToken.Valid {
        return 0, time.Time{}, jwt.ErrSignatureInvalid
    }

    // extract claims
    claims, ok := parsedToken.Claims.(jwt.MapClaims)
    if !ok {
        return 0, time.Time{}, jwt.ErrInvalidKey
    }

    accountID, ok := claims["sub"].(float64)
    if !ok {
        return 0, time.Time{}, jwt.ErrInvalidKey
    }

    expireTimeUnix, ok := claims["exp"].(float64)
    if !ok {
        return 0, time.Time{}, jwt.ErrInvalidKey
    }

    expireTime := time.Unix(int64(expireTimeUnix), 0)
    return uint64(accountID), expireTime, nil
}

// GetToken implements TokenHandler.
func (t tokenHandler) GetToken(ctx context.Context, accountID uint64) (string, time.Time, error) {
	expireTime := time.Now().Add(t.expiresIn)
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"kid": t.publicKeyID,
		"sub": accountID,
		"exp": expireTime,
	})

	signedToken, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", time.Time{}, err

	}
	return signedToken, expireTime, nil
}

func (t tokenHandler) WithDatabase(database database.Database) TokenHandler {
    t.accountDataAccessor = t.accountDataAccessor.WithDatabase(database)
    return t
}

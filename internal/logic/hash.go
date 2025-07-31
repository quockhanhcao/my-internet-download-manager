package logic

import (
	"context"
	"errors"
	"log"

	"github.com/quockhanhcao/my-internet-download-manager/internal/configs"
	"golang.org/x/crypto/bcrypt"
)

type Hash interface {
	Hash(ctx context.Context, payload string) (string, error)
	IsHashEqual(ctx context.Context, hash, data string) (bool, error)
}

type hash struct {
	accountConfig configs.AccountConfig
}

func NewHash(accountConfig configs.AccountConfig) Hash {
	return &hash{
		accountConfig: accountConfig,
	}
}

func (h *hash) Hash(ctx context.Context, payload string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(payload), h.accountConfig.HashCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		return "", err
	}
	return string(result), nil
}

func (h *hash) IsHashEqual(ctx context.Context, hash string, data string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		log.Printf("error comparing hash: %v", err)
		return false, err
	}
	return true, nil
}

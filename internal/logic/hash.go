package logic

import (
	"context"
	"errors"
	"log"

	"github.com/quockhanhcao/my-internet-download-manager/internal/configs"
	"golang.org/x/crypto/bcrypt"
)

type HashHandler interface {
	HashPassword(ctx context.Context, password string) (string, error)
	IsHashEqual(ctx context.Context, hashedPassword, inputPassword string) (bool, error)
}

type hash struct {
	accountConfigs configs.AccountConfig
}

func NewHashHandler(accountConfigs configs.AccountConfig) HashHandler {
	return &hash{
		accountConfigs: accountConfigs,
	}
}

func (h hash) HashPassword(ctx context.Context, password string) (string, error) {
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), h.accountConfigs.HashCost)
	if err != nil {
		log.Print("error hashing password:", err)
		return "", err
	}
	return string(bcryptHash), nil
}

func (h hash) IsHashEqual(ctx context.Context, hash, data string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Print("passwords do not match:", err)
			return false, nil
		}
		log.Print("error comparing hashed password", err)
		return false, err
	}
	return true, nil
}

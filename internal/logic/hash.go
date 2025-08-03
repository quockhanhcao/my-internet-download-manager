package logic

import (
	"errors"
	"log"

	"github.com/quockhanhcao/my-internet-download-manager/internal/configs"
	"golang.org/x/crypto/bcrypt"
)

type Hash interface {
	HashPassword(password string) (string, error)
	IsHashEqual(hashedPassword, password string) (bool, error)
}

type hash struct {
	accountConfigs configs.AccountConfig
}

func NewHash(accountConfigs configs.AccountConfig) Hash {
	return &hash{
		accountConfigs: accountConfigs,
	}
}

func (h hash) HashPassword(password string) (string, error) {
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), h.accountConfigs.HashCost)
	if err != nil {
		log.Print("error hashing password:", err)
		return "", err
	}
	return string(bcryptHash), nil
}

func (h hash) IsHashEqual(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
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

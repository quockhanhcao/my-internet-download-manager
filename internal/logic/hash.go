package logic

import (
	"context"
	"errors"

	"github.com/quockhanhcao/my-internet-download-manager/internal/configs"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type HashHandler interface {
	Hash(ctx context.Context, password string) (string, error)
	IsHashEqual(ctx context.Context, hashedPassword, inputPassword string) (bool, error)
}

type hash struct {
	configs configs.AuthConfig
	logger  *zap.Logger
}

func NewHashHandler(configs configs.AuthConfig, logger *zap.Logger) HashHandler {
	return &hash{
		configs: configs,
		logger:  logger,
	}
}

func (h hash) Hash(ctx context.Context, password string) (string, error) {
	h.logger.With(zap.Int("passwordLength", len(password)), zap.Int("hashCost", h.configs.HashConfig.HashCost)).Info("hashing password")

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), h.configs.HashConfig.HashCost)
	if err != nil {
		h.logger.With(zap.Error(err), zap.Int("hashCost", h.configs.HashConfig.HashCost)).Error("failed to hash password")
		return "", err
	}

	h.logger.With(zap.Int("passwordLength", len(password)), zap.Int("hashLength", len(bcryptHash))).Info("password hashed successfully")
	return string(bcryptHash), nil
}

func (h hash) IsHashEqual(ctx context.Context, hash, data string) (bool, error) {
	h.logger.With(zap.Int("hashLength", len(hash)), zap.Int("dataLength", len(data))).Info("comparing hash and password")

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			h.logger.With(zap.Int("hashLength", len(hash)), zap.Int("dataLength", len(data))).Warn("password does not match hash")
			return false, nil
		}
		h.logger.With(zap.Error(err), zap.Int("hashLength", len(hash))).Error("failed to compare hash and password")
		return false, err
	}

	h.logger.With(zap.Int("hashLength", len(hash)), zap.Int("dataLength", len(data))).Info("password matches hash successfully")
	return true, nil
}

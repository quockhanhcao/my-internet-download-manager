package configs

import "time"

type AuthConfig struct {
	KeyBitSize int           `yaml:"key_bit_size"`
	ExpireTime time.Duration `yaml:"expire_time"`
}

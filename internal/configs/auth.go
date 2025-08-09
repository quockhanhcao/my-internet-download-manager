package configs

import "time"

type HashConfig struct {
	HashCost int `yaml:"hash_cost"`
}

type TokenConfig struct {
	KeyBitSize                  int    `yaml:"key_bit_size"`
	ExpiresIn                   string `yaml:"expires_in"`
	RegenerateTokenBeforeExpiry string `yaml:"regenerate_token_before_expiry"`
}

type AuthConfig struct {
	HashConfig  HashConfig
	TokenConfig TokenConfig
}

func (t TokenConfig) GetExpiresInDuration() (time.Duration, error) {
	return time.ParseDuration(t.ExpiresIn)
}

func (t TokenConfig) GetRegenerateTokenBeforeExpiryDuration() (time.Duration, error) {
	return time.ParseDuration(t.RegenerateTokenBeforeExpiry)
}

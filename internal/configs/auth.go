package configs

import "time"

type AuthConfig struct {
	ExpireTime time.Duration `yaml:"expire_time"`
}

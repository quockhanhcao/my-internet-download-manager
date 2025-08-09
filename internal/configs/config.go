package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigFilePath string

type Config struct {
	DatabaseConfig DatabaseConfig `yaml:"database_config"`
	AuthConfig     AuthConfig     `yaml:"auth_config"`
}

func NewConfig(filePath ConfigFilePath) (Config, error) {
	configBytes, err := os.ReadFile(string(filePath))
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

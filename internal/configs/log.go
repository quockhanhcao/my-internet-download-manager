package configs

type LogConfig struct {
	Level       string   `mapstructure:"level"`
	OutputPaths []string `mapstructure:"output_paths"`
}

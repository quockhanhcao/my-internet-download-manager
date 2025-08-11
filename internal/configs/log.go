package configs

type LogConfig struct {
	Level       string   `yaml:"level"`
	OutputPaths []string `yaml:"output_paths"`
}

package config

import (
	"github.com/BurntSushi/toml"
	"monkey-test-api/internal/logger"
)

type Config struct {
	Log      logger.Config `toml:"log"`
	Database struct {
		Type string `toml:"type"`
		MySQL struct {
			DSN string `toml:"dsn"`
		} `toml:"mysql"`
		MongoDB struct {
			URI      string `toml:"uri"`
			Database string `toml:"database"`
		} `toml:"mongodb"`
	} `toml:"database"`
	// ... 其他配置
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

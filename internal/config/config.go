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
	Task struct {
		String ParamTestConfig `toml:"string"`
		Number ParamTestConfig `toml:"number"`
		Bool   ParamTestConfig `toml:"bool"`
	} `toml:"task"`
}

type ParamTestConfig struct {
	Values []ParamTestValue `toml:"values"`
}

type ParamTestValue struct {
	Value  interface{} `toml:"value"`
	Reason string      `toml:"reason"`
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type DatabaseConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	DbName   string `toml:"dbname"`
}

type Config struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &config, nil
}

package main

import (
	"fmt"
	"log"
	"monkey-test-api/internal/config"
)

func main() {
	// 假设配置文件路径为 ./config.toml
	configFile := "./conf/config.toml"

	// 加载配置
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// 输出配置内容
	fmt.Println("Server Config:", cfg.Server)
	fmt.Println("Database Config:", cfg.Database)

}

package config

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Version 版本号
const Version = "0.0.1"

const cfgPath = "config.toml"

// ReadConf 读取配置
func ReadConf() Config {
	initCfg()
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v\n", err)
	}
	var config Config
	err = toml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling config file: %v\n", err)
	}
	return config
}

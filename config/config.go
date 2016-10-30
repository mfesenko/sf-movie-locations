package config

import (
	"os"

	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server     ServerConfig
	Db         DbConfig
	Dataloader DataloaderConfig
}

type DbConfig struct {
	Host           string
	DbName         string
	CollectionName string
	Username       string
	Password       string
	Timeout        int
}

type ServerConfig struct {
	Port              int
	StaticContentPath string
}

type DataloaderConfig struct {
	BatchSize int
}

func LoadConfigFile(filePath string) (Config, error) {
	var config Config
	_, err := os.Stat(filePath)
	if err != nil {
		return config, fmt.Errorf("Config file is not present at path %s", filePath)
	}
	_, err = toml.DecodeFile(filePath, &config)
	if err != nil {
		return config, fmt.Errorf("Failed to parse config file: %s", err)
	}
	return config, nil
}

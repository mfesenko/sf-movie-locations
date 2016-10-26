package server

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server ServerConfig
	Db     DbConfig
}

type ServerConfig struct {
	Port int
}

type DbConfig struct {
	Host           string
	DbName         string
	CollectionName string
}

func LoadConfigFile(filePath string) Config {
	_, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("Config file is not present at path %s", filePath)
	}
	var config Config
	_, err = toml.DecodeFile(filePath, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err)
	}
	return config
}

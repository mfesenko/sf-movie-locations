package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDataDir = "testdata/"

func TestLoadNotExistingConfigFile(t *testing.T) {
	filePath := testDataDir + "manticora"
	config, err := LoadConfigFile(filePath)
	assert := assert.New(t)
	assert.EqualError(err, "Config file is not present at path "+filePath)
	assert.Equal(Config{}, config)
}

func TestLoadNotTomlConfigFile(t *testing.T) {
	config, err := LoadConfigFile(testDataDir + "not-toml-config")
	assert := assert.New(t)
	assert.EqualError(err, "Failed to parse config file: Near line 1 (last key parsed 'hello'): Expected key separator '=', but got 'w' instead.")
	assert.Equal(Config{}, config)
}

func TestLoadTomlConfigFileWithInvalidType(t *testing.T) {
	config, err := LoadConfigFile(testDataDir + "invalid-toml-config.toml")
	assert := assert.New(t)
	assert.EqualError(err, "Failed to parse config file: Near line 2 (last key parsed 'dataloader.batchSize'): Expected value but found \"asdf\" instead.")
	assert.Equal(Config{}, config)
}

func TestLoadConfigFile(t *testing.T) {
	config, err := LoadConfigFile(testDataDir + "valid-toml-config.toml")
	assert := assert.New(t)
	assert.Nil(err)
	expectedConfig := Config{
		Dataloader: DataloaderConfig{
			BatchSize: 100,
		},
		Server: ServerConfig{
			Port:              8080,
			StaticContentPath: "static",
		},
		Db: DbConfig{
			Host:           "localhost",
			DbName:         "db-name",
			CollectionName: "collection-name",
			Username:       "username",
			Password:       "pwd",
			Timeout:        30,
		},
	}
	assert.Equal(expectedConfig, config)
}

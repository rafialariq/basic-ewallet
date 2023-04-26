package config

import "final_project_easycash/utils"

type ApiConfig struct {
	ServerPort string
}

type DbConfig struct {
	Host, Port, User, Password, Name, SslMode string
}

type StorageConfig struct {
	BaseFilePath string
}

type AppConfig struct {
	ApiConfig
	DbConfig
	StorageConfig
}

func (c *AppConfig) readConfigFile() {
	c.DbConfig = DbConfig{
		Host:     utils.DotEnv("DB_HOST"),
		Port:     utils.DotEnv("DB_PORT"),
		User:     utils.DotEnv("DB_USER"),
		Password: utils.DotEnv("DB_PASSWORD"),
		Name:     utils.DotEnv("DB_NAME"),
		SslMode:  utils.DotEnv("SSL_MODE"),
	}
	c.ApiConfig = ApiConfig{
		ServerPort: utils.DotEnv("SERVER_PORT"),
	}
	c.StorageConfig = StorageConfig{
		BaseFilePath: utils.DotEnv("BASE_FILE_PATH"),
	}
}

func NewConfig() AppConfig {
	config := AppConfig{}
	config.readConfigFile()
	return config
}

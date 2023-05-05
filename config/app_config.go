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
	envFilePath := ".env"
	c.DbConfig = DbConfig{
		Host:     utils.DotEnv("DB_HOST", envFilePath),
		Port:     utils.DotEnv("DB_PORT", envFilePath),
		User:     utils.DotEnv("DB_USER", envFilePath),
		Password: utils.DotEnv("DB_PASSWORD", envFilePath),
		Name:     utils.DotEnv("DB_NAME", envFilePath),
		SslMode:  utils.DotEnv("SSL_MODE", envFilePath),
	}
	c.ApiConfig = ApiConfig{
		ServerPort: utils.DotEnv("SERVER_PORT", envFilePath),
	}
	c.StorageConfig = StorageConfig{
		BaseFilePath: utils.DotEnv("BASE_FILE_PATH", envFilePath),
	}
}

func NewConfig() AppConfig {
	config := AppConfig{}
	config.readConfigFile()
	return config
}

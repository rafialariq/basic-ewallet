package manager

import (
	"final_project_easycash/config"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type InfraManager interface {
	ConnectDb() *sqlx.DB
	InitializeBasePath() string
}

type infraManager struct {
	db     *sqlx.DB
	config config.AppConfig
}

func (i *infraManager) initDb() {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", i.config.User, i.config.Password, i.config.Host, i.config.Port, i.config.Name, i.config.SslMode)
	db, err := sqlx.Connect("postgres", dataSourceName)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("Application filed to run", err)
			db.Close()
		}
	}()

	i.db = db
	fmt.Println("Connected to DB")
}

func (i *infraManager) ConnectDb() *sqlx.DB {
	return i.db
}

func (i *infraManager) InitializeBasePath() string {
	return i.config.BaseFilePath
}

func NewInfraManager(config config.AppConfig) InfraManager {
	infra := infraManager{
		config: config,
	}
	infra.initDb()
	return &infra
}

package main

import (
	"final_project_easycash/controller"
	"final_project_easycash/repository"
	"final_project_easycash/usecase"
	"final_project_easycash/utils"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := utils.DotEnv("DB_HOST")
	dbPort := utils.DotEnv("DB_PORT")
	dbUser := utils.DotEnv("DB_USER")
	dbPassword := utils.DotEnv("DB_PASSWORD")
	dbName := utils.DotEnv("DB_NAME")
	sslMode := utils.DotEnv("SSL_MODE")
	serverPort := utils.DotEnv("SERVER_PORT")

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)
	db, err := sqlx.Connect("postgres", dataSourceName)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Println("Failed to close DB connection")
		}
	}()

	if err != nil {
		log.Println("Failed to connect to DB")
		os.Exit(1)
	} else {
		log.Println("Connected to DB")
	}

	transactionRepo := repository.NewTransactionRepo(db)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo)
	transactionController := controller.NewTransactionController(transactionUsecase)

	router := gin.Default()

	router.POST("/merchant", transactionController.TransferMoney)

	err = router.Run(serverPort)
	if err != nil {
		log.Println("Failed to start server")
		os.Exit(1)
	}
}

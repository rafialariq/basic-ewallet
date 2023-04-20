package main

import (
	"final_project_easycash/controller"
	"final_project_easycash/repository"
	"final_project_easycash/usecase"
	"final_project_easycash/utils"
	"fmt"
	"log"

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
	baseFilePath := utils.DotEnv("BASE_FILE_PATH")

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)
	db, err := sqlx.Connect("postgres", dataSourceName)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to DB")
	}

	fileRepo := repository.NewFileRepository(baseFilePath)
	userRepo := repository.NewUserRepo(db)
	userUsecase := usecase.NewUserUsecase(userRepo, fileRepo)
	userController := controller.NewUserController(userUsecase)

	router := gin.Default()

	userProfileRouter := router.Group("/profile")
	userProfileRouter.GET("/:id", userController.CheckProfile)
	userProfileRouter.POST("/edit", userController.EditProfile)
	userProfileRouter.POST("/edit/photo/:id", userController.EditPhotoProfile)
	userProfileRouter.DELETE("/:id", userController.UnregProfile)

	if err := router.Run(serverPort); err != nil {
		log.Fatal(err)
	}
}

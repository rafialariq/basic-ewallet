package main

import (
	"final_project_easycash/controller"
	"final_project_easycash/middleware"
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
	transferRepo := repository.NewTransferRepo(db)
	transferUsecase := usecase.NewTransferUsecase(transferRepo)
	transferController := controller.NewTransferController(transferUsecase, userUsecase)
	topUpRepo := repository.NewTopUpRepo(db)
	topUpUsecase := usecase.NewTopUpUsecase(topUpRepo)
	topUpController := controller.NewTopUpController(topUpUsecase, userUsecase)
	registerRepo := repository.NewRegisterRepo(db)
	registerService := usecase.NewRegisterService(registerRepo)
	registerController := controller.NewRegisterController(registerService)

	loginRepo := repository.NewLoginRepo(db)
	loginService := usecase.NewLoginService(loginRepo)
	loginController := controller.NewLoginController(loginService)

	router := gin.Default()

	menuRouter := router.Group("/menu")
	menuRouter.Use(middleware.AuthMiddleware())

	menuRouter.GET("/profile/:username", userController.CheckProfile)
	menuRouter.POST("/profile/edit", userController.EditProfile)
	menuRouter.POST("/profile/edit/photo/:username", userController.EditPhotoProfile)
	menuRouter.DELETE("/profile/:username", userController.UnregProfile)
	menuRouter.POST("/transfer/bank", transferController.WithdrawBalance)
	menuRouter.POST("/transfer/user", transferController.TransferBalance)
	menuRouter.POST("/topup", topUpController.TopUpBalance)
	router.POST("/signup", registerController.RegisterHandler)
	router.POST("/login", loginController.LoginHandler)

	if err := router.Run(serverPort); err != nil {
		log.Fatal(err)
	}
}

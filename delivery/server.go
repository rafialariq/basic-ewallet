package delivery

import (
	"final_project_easycash/config"
	"final_project_easycash/controller"
	"final_project_easycash/manager"
	"final_project_easycash/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

type AppServer struct {
	usecaseManager manager.UsecaseManager
	engine         *gin.Engine
	host           string
}

func (p *AppServer) menu() {
	menuRoutes := p.engine.Group("/menu")
	menuRoutes.Use(middleware.AuthMiddleware())
	p.userController(menuRoutes)
	p.transferController(menuRoutes)
	p.topUpController(menuRoutes)
	p.registerController(p.engine)
	p.loginController(p.engine)
}

func (p *AppServer) userController(rg *gin.RouterGroup) {
	controller.NewUserController(rg, p.usecaseManager.UserUsecase())
}

func (p *AppServer) transferController(rg *gin.RouterGroup) {
	controller.NewTransferController(rg, p.usecaseManager.TransferUsecase(), p.usecaseManager.UserUsecase())
}

func (p *AppServer) topUpController(rg *gin.RouterGroup) {
	controller.NewTopUpController(rg, p.usecaseManager.TopUpUsecase(), p.usecaseManager.UserUsecase())
}

func (p *AppServer) registerController(r *gin.Engine) {
	controller.NewRegisterController(r, p.usecaseManager.RegisterUsecase())
}

func (p *AppServer) loginController(r *gin.Engine) {
	controller.NewLoginController(r, p.usecaseManager.LoginUsecase())
}

func (p *AppServer) Run() {
	p.menu()
	err := p.engine.Run(p.host)
	defer func() {
		if err := recover(); err != nil {
			log.Println("Application failed to run", err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
}

func Server() *AppServer {
	router := gin.Default()
	config := config.NewConfig()
	infraManager := manager.NewInfraManager(config)
	repoManager := manager.NewRepoManager(infraManager)
	usecaseManager := manager.NewUsecaseManager(repoManager)
	host := config.ServerPort
	return &AppServer{
		usecaseManager: usecaseManager,
		engine:         router,
		host:           host,
	}
}

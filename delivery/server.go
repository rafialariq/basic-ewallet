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
	routes := p.engine.Group("/")
	routes.Use(middleware.LoggingMiddleware("../.log"))
	menuRoutes := routes.Group("/menu")
	menuRoutes.Use(middleware.AuthMiddleware())
	p.userController(menuRoutes)
	p.transactionController(menuRoutes)
	p.registerController(routes)
	p.loginController(routes)
	p.historyController(menuRoutes)
}

func (p *AppServer) userController(r *gin.RouterGroup) {
	controller.NewUserController(r, p.usecaseManager.UserUsecase())
}

func (p *AppServer) transactionController(rg *gin.RouterGroup) {
	controller.NewTransactionController(rg, p.usecaseManager.TransactionUsecase(), p.usecaseManager.UserUsecase())
}

func (p *AppServer) registerController(r *gin.RouterGroup) {
	controller.NewRegisterController(r, p.usecaseManager.RegisterUsecase())
}

func (p *AppServer) loginController(r *gin.RouterGroup) {
	controller.NewLoginController(r, p.usecaseManager.LoginUsecase())
}

func (p *AppServer) historyController(rg *gin.RouterGroup) {
	controller.NewHistoryController(rg, p.usecaseManager.HistoryUsecase())
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

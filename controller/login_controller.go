package controller

import (
	"net/http"

	"final_project_easycash/model"
	"final_project_easycash/usecase"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	loginService usecase.LoginService
}

func (l *LoginController) LoginHandler(ctx *gin.Context) {
	var user model.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recUser, res := l.loginService.UserLogin(user)

	if recUser {
		ctx.JSON(http.StatusOK, gin.H{"token": res})
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": res})
	}

}

func NewLoginController(r *gin.Engine, u usecase.LoginService) *LoginController {
	controller := LoginController{
		loginService: u,
	}
	r.POST("/login", controller.LoginHandler)
	return &controller
}

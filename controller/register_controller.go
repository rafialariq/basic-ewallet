package controller

import (
	"fmt"
	"net/http"

	"final_project_easycash/model"
	"final_project_easycash/usecase"

	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	registerService usecase.RegisterService
}

func (r *RegisterController) RegisterHandler(ctx *gin.Context) {
	var newUser model.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(newUser)
	user, res := r.registerService.UserSignup(&newUser)

	if user {
		ctx.JSON(http.StatusCreated, gin.H{
			"msg":   "user created successfully",
			"token": res,
		})
		return
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": res})
	}

}

func NewRegisterController(u usecase.RegisterService) *RegisterController {
	controller := RegisterController{
		registerService: u,
	}

	return &controller
}

package controller

import (
	"final_project_easycash/model"
	"final_project_easycash/usecase"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TopUpController struct {
	usecase     usecase.TopUpUsecase
	usecaseUser usecase.UserUsecase
}

func (c *TopUpController) TopUpBalance(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "missing claims"})
		return
	}

	usernameToken, ok := claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	userToken, err := c.usecaseUser.CheckProfile(usernameToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bill model.Bill

	if err := ctx.ShouldBind(&bill); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userToken.PhoneNumber != bill.DestinationId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res := c.usecase.TopUpBalance(bill.SenderId, bill.DestinationId, bill.Amount)

	if res != nil {
		if res.Error() == "Receiver number not found" || res.Error() == "Sender number not found" || res.Error() == "Balance is not sufficient" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": res.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error()})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "transaction added"})
}

func NewTopUpController(u usecase.TopUpUsecase, us usecase.UserUsecase) *TopUpController {
	controller := TopUpController{
		usecase:     u,
		usecaseUser: us,
	}
	return &controller
}

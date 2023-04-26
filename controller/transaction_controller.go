package controller

import (
	"errors"
	"final_project_easycash/model"
	"final_project_easycash/usecase"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	usecase     usecase.TransactionUsecase
	usecaseUser usecase.UserUsecase
}

func (c *TransactionController) TransferMoney(ctx *gin.Context) {
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
		ctx.JSON(http.StatusBadRequest, "invalid input")
		return
	}

	if userToken.PhoneNumber != bill.SenderId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if bill.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, errors.New("invalid amount"))
		return
	}

	err = c.usecase.TransferMoney(bill.SenderId, bill.DestinationId, bill.Amount)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, "transaction added")
}

func (c *TransactionController) TopUpBalance(ctx *gin.Context) {
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

func (c *TransactionController) WithdrawBalance(ctx *gin.Context) {
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

	if userToken.PhoneNumber != bill.SenderId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res := c.usecase.WithdrawBalance(bill.SenderId, bill.DestinationId, bill.Amount)

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

func (c *TransactionController) TransferBalance(ctx *gin.Context) {
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

	if userToken.PhoneNumber != bill.SenderId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res := c.usecase.TransferBalance(bill.SenderId, bill.DestinationId, bill.Amount)

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

func NewTransactionController(rg *gin.RouterGroup, u usecase.TransactionUsecase, us usecase.UserUsecase) *TransactionController {
	controller := TransactionController{
		usecase:     u,
		usecaseUser: us,
	}
	rg.POST("/topup", controller.TopUpBalance)
	rg.POST("/transfer/bank", controller.WithdrawBalance)
	rg.POST("/transfer/user", controller.TransferBalance)
	rg.POST("/merchant", controller.TransferMoney)
	return &controller
}

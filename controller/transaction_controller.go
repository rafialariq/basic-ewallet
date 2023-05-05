package controller

import (
	"encoding/json"
	"errors"
	"final_project_easycash/model"
	"final_project_easycash/usecase"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	usecase     usecase.TransactionUsecase
	usecaseUser usecase.UserUsecase
}

func (c *TransactionController) TransferMoney(ctx *gin.Context) {
	var bill model.Bill

	if err := ctx.ShouldBind(&bill); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid input")
		return
	}

	if bill.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, errors.New("invalid amount"))
		return
	}

	err := c.usecase.TransferMoney(bill.SenderId, bill.DestinationId, bill.Amount)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, "transaction added")
}

func (c *TransactionController) TopUpBalance(ctx *gin.Context) {
	var bill model.Bill

	if err := ctx.ShouldBindJSON(&bill); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := c.usecase.TopUpBalance(bill.SenderId, bill.DestinationId, bill.Amount)

	if res != nil {
		if res.Error() == "Receiver number not found" || res.Error() == "Sender number not found" || res.Error() == "Balance is not sufficient" || res.Error() == "Minimum Transaction Rp 10.000,00" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": res.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "transaction added"})
}

func (c *TransactionController) WithdrawBalance(ctx *gin.Context) {
	rawBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var bill model.Bill
	if err := json.Unmarshal(rawBody, &bill); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
		return
	}

	usernameToken, ok := claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	userToken, err := c.usecaseUser.CheckProfile(usernameToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userToken.PhoneNumber != bill.SenderId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res := c.usecase.WithdrawBalance(bill.SenderId, bill.DestinationId, bill.Amount)

	if res != nil {
		if res.Error() == "Receiver number not found" || res.Error() == "Sender number not found" || res.Error() == "Balance is not sufficient" || res.Error() == "Minimum Transaction Rp 10.000,00" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": res.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "transaction added"})
}

func (c *TransactionController) TransferBalance(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
		return
	}

	usernameToken, ok := claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		log.Print("cek")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	userToken, err := c.usecaseUser.CheckProfile(usernameToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var bill model.Bill
	if err := json.Unmarshal(rawBody, &bill); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userToken.PhoneNumber != bill.SenderId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res := c.usecase.TransferBalance(bill.SenderId, bill.DestinationId, bill.Amount)

	log.Print(res)

	if res != nil {
		if res.Error() == "Receiver number not found" || res.Error() == "Sender number not found" || res.Error() == "Balance is not sufficient" || res.Error() == "Minimum Transaction Rp 10.000,00" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": res.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "transaction added"})
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

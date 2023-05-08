package controller

import (
	"encoding/json"
	"errors"
	"final_project_easycash/model"
	"final_project_easycash/repository"
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

	if err := ctx.ShouldBindJSON(&bill); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if bill.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	err := c.usecase.TransferMoney(bill.SenderId, bill.DestinationId, bill.Amount)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "transaction added"})
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

func (c *TransactionController) SplitBill(ctx *gin.Context) {
	var req struct {
		Sender   string    `json:"sender_id"`
		Receiver []string  `json:"destination_id"`
		Amount   []float64 `json:"amount"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Receiver) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "must provide at least one receiver"})
		return
	}

	totalAmount := 0.0
	for _, amount := range req.Amount {
		if amount <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
			return
		}
		totalAmount += amount
	}

	err := c.usecase.SplitBill(req.Sender, req.Receiver, req.Amount)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "bill split successfully"})
}

func (c *TransactionController) PayBill(ctx *gin.Context) {
	id_transaction := ctx.PostForm("idTransaction")
	receiver := ctx.PostForm("receiver")

	if id_transaction == "" || receiver == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := c.usecase.PayBill(receiver, id_transaction)
	if err != nil {
		if errors.Is(err, repository.ErrBillNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
			return
		} else if errors.Is(err, repository.ErrBillPaid) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Bill has already been paid"})
			return
		} else if errors.Is(err, repository.ErrInsufficientBalance) {
			ctx.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient balance"})
			return
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
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
	rg.POST("/split-bill", controller.SplitBill)
	rg.POST("/PayBill", controller.PayBill)
	return &controller
}

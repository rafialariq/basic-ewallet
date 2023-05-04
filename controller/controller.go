package controller

import (
	"final_project_easycash/model"
	"final_project_easycash/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	usecase usecase.TransactionUsecase
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

func NewTransactionController(u usecase.TransactionUsecase) *TransactionController {
	controller := TransactionController{
		usecase: u,
	}
	return &controller
}

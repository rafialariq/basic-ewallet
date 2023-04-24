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

func (c *TransactionController) WithdrawBalance(ctx *gin.Context) {
	var bill model.Bill

	if err := ctx.ShouldBind(&bill); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := c.usecase.WithdrawBalance(bill.SenderId, bill.DestinationId, bill.Amount)

	if res != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "transaction added"})
}

func (c *TransactionController) TransferBalance(ctx *gin.Context) {
	var bill model.Bill

	if err := ctx.ShouldBind(&bill); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := c.usecase.TransferBalance(bill.SenderId, bill.DestinationId, bill.Amount)

	if res != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "transaction added"})
}

func NewTransactionController(u usecase.TransactionUsecase) *TransactionController {
	controller := TransactionController{
		usecase: u,
	}
	return &controller
}

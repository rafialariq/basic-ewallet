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
		ctx.JSON(http.StatusBadRequest, "invalid input")
		return
	}

	res := c.usecase.WithdrawBalance(bill.SenderId, bill.DestinationId, bill.Amount)

	if res != nil {
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	ctx.JSON(http.StatusCreated, "transaction added")
}

func NewTransactionController(u usecase.TransactionUsecase) *TransactionController {
	controller := TransactionController{
		usecase: u,
	}
	return &controller
}

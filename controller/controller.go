package controller

import (
	"errors"
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

func NewTransactionController(u usecase.TransactionUsecase) *TransactionController {
	controller := TransactionController{
		usecase: u,
	}
	return &controller
}

package controller

import (
	"errors"
	"final_project_easycash/model"
	"final_project_easycash/repository"
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

/*func (c *TransactionController) PayBill(w http.ResponseWriter, r *http.Request) {
	idTransaction := r.FormValue("idTransaction")
	receiver := r.FormValue("receiver")

	if idTransaction == "" || receiver == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := c.usecase.PayBill(receiver, idTransaction)
	if err != nil {
		if errors.Is(err, repository.ErrBillNotFound) {
			http.Error(w, "Bill not found", http.StatusNotFound)
			return
		} else if errors.Is(err, repository.ErrBillPaid) {
			http.Error(w, "Bill has already been paid", http.StatusForbidden)
			return
		} else if errors.Is(err, repository.ErrInsufficientBalance) {
			http.Error(w, "Insufficient balance", http.StatusPaymentRequired)
			return
		} else {
			http.Error(w, "Failed to process payment", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment processed successfully"))
}*/

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

func NewTransactionController(u usecase.TransactionUsecase) *TransactionController {
	controller := TransactionController{
		usecase: u,
	}
	return &controller
}

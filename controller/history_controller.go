package controller

import (
	"net/http"
	"strconv"

	"final_project_easycash/model"
	"final_project_easycash/usecase"

	"github.com/gin-gonic/gin"
)

type HistoryController struct {
	historyUsecase usecase.HistoryUsecase
}

func (h *HistoryController) FindAllByUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.historyUsecase.FindAllByUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)

}

func (h *HistoryController) FindByAccountType(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountTypeId, err := strconv.Atoi(ctx.Param("accountTypeId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.historyUsecase.FindByAccountType(user, accountTypeId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)

}

func (h *HistoryController) FindByType(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	typeId := ctx.Param("typeId")

	res, err := h.historyUsecase.FindByType(user, typeId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *HistoryController) FindByAmount(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	moreThan, err := strconv.ParseFloat(ctx.Param("more_than"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lessThan, err := strconv.ParseFloat(ctx.Param("less_than"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.historyUsecase.FindByAmount(user, moreThan, lessThan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func NewHistoryController(u usecase.HistoryUsecase) *HistoryController {
	controller := HistoryController{
		historyUsecase: u,
	}

	return &controller
}

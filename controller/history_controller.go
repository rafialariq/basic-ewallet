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

	res, err := h.historyUsecase.HistoryByUser(user)
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

	res, err := h.historyUsecase.HistoryWithAccountFilter(user, accountTypeId)
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

	typeId, err := strconv.Atoi(ctx.Param("typeId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.historyUsecase.HistoryWithTypeFilter(user, typeId)
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

	res, err := h.historyUsecase.HistoryWithAmountFilter(user, moreThan, lessThan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func NewHistoryController(rg *gin.RouterGroup, u usecase.HistoryUsecase) *HistoryController {
	controller := HistoryController{
		historyUsecase: u,
	}
	rg.GET("/history", controller.FindAllByUser)
	rg.GET("/history/account/:accountTypeId", controller.FindByAccountType)
	rg.GET("/history/type/:typeId", controller.FindByType)
	rg.GET("/history/:more_than/:less_than", controller.FindByAmount)
	return &controller
}

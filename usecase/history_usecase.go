package usecase

import (
	"final_project_easycash/model"
	"final_project_easycash/repository"
)

type HistoryUsecase interface {
	HistoryByUser(user model.User) ([]model.Bill, error)
	HistoryWithAccountFilter(user model.User, accountTypeId int) ([]model.Bill, error)
	HistoryWithTypeFilter(user model.User, typeId int) ([]model.Bill, error)
	HistoryWithAmountFilter(user model.User, moreThan, lessThan float64) ([]model.Bill, error)
}

type historyUsecase struct {
	historyRepo repository.HistoryRepo
}

func (h *historyUsecase) HistoryByUser(user model.User) ([]model.Bill, error) {
	return h.historyRepo.GetHistoryByUser(user)
}

func (h *historyUsecase) HistoryWithAccountFilter(user model.User, accountTypeId int) ([]model.Bill, error) {
	return h.historyRepo.GetHistoryWithAccountFilter(user, accountTypeId)
}

func (h *historyUsecase) HistoryWithTypeFilter(user model.User, typeId int) ([]model.Bill, error) {
	return h.historyRepo.GetHistoryWithTypeFilter(user, typeId)
}

func (h *historyUsecase) HistoryWithAmountFilter(user model.User, moreThan, lessThan float64) ([]model.Bill, error) {
	return h.historyRepo.GetHistoryWithAmountFilter(user, moreThan, lessThan)
}

func NewHistoryUsecase(historyRepo repository.HistoryRepo) HistoryUsecase {
	return &historyUsecase{
		historyRepo: historyRepo,
	}
}

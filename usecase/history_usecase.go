package usecase

import (
	"final_project_easycash/model"
	"final_project_easycash/repository"
)

type HistoryUsecase interface {
	FindAllByUser(user model.User) ([]model.Bill, error)
	FindByAccountType(user model.User, accountTypeId int) ([]model.Bill, error)
	FindByType(user model.User, typeId string) ([]model.Bill, error)
	FindByAmount(user model.User, moreThan, lessThan float64) ([]model.Bill, error)
}

type historyUsecase struct {
	historyRepo repository.HistoryRepo
}

func (h *historyUsecase) FindAllByUser(user model.User) ([]model.Bill, error) {
	return h.historyRepo.GetAllByUser(user)
}

func (h *historyUsecase) FindByAccountType(user model.User, accountTypeId int) ([]model.Bill, error) {
	return h.historyRepo.GetByAccountType(user, accountTypeId)
}

func (h *historyUsecase) FindByType(user model.User, typeId string) ([]model.Bill, error) {
	return h.historyRepo.GetByType(user, typeId)
}

func (h *historyUsecase) FindByAmount(user model.User, moreThan, lessThan float64) ([]model.Bill, error) {
	return h.historyRepo.GetByAmount(user, moreThan, lessThan)
}

func NewHistoryUsecase(historyRepo repository.HistoryRepo) HistoryUsecase {
	return &historyUsecase{
		historyRepo: historyRepo,
	}
}

package usecase

import (
	"final_project_easycash/repository"
)

type TopUpUsecase interface {
	TopUpBalance(sender string, receiver string, amount float64) error
}

type topUpUsecase struct {
	topUpRepo repository.TopUpRepo
}

func (u *topUpUsecase) TopUpBalance(sender string, receiver string, amount float64) error {
	return u.topUpRepo.TopUpBalance(sender, receiver, amount)
}

func NewTopUpUsecase(topUpRepo repository.TopUpRepo) TopUpUsecase {
	return &topUpUsecase{
		topUpRepo: topUpRepo,
	}
}

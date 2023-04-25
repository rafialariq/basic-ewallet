package usecase

import (
	"final_project_easycash/repository"
)

type TransferUsecase interface {
	WithdrawBalance(sender string, receiver string, amount float64) error
	TransferBalance(sender string, receiver string, amount float64) error
}

type transferUsecase struct {
	transferRepo repository.TransferRepo
}

func (u *transferUsecase) WithdrawBalance(sender string, receiver string, amount float64) error {
	return u.transferRepo.WithdrawBalance(sender, receiver, amount)
}

func (u *transferUsecase) TransferBalance(sender string, receiver string, amount float64) error {
	return u.transferRepo.TransferBalance(sender, receiver, amount)
}

func NewTransferUsecase(transferRepo repository.TransferRepo) TransferUsecase {
	return &transferUsecase{
		transferRepo: transferRepo,
	}
}

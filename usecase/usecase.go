package usecase

import (
	"final_project_easycash/repository"
)

type TransactionUsecase interface {
	TransferMoney(sender string, receiver string, amount float64) error
}

type transactionUsecase struct {
	transactionRepo repository.TransactionRepo
}

func (u *transactionUsecase) TransferMoney(sender string, receiver string, amount float64) error {
	return u.transactionRepo.TransferMoney(sender, receiver, amount)
}

func NewTransactionUsecase(transactionRepo repository.TransactionRepo) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
	}
}

package usecase

import (
	"final_project_easycash/repository"
)

type TransactionUsecase interface {
	TransferMoney(sender string, receiver string, amount float64) error
	TopUpBalance(sender string, receiver string, amount float64) error
	WithdrawBalance(sender string, receiver string, amount float64) error
	TransferBalance(sender string, receiver string, amount float64) error
}

type transactionUsecase struct {
	transactionRepo repository.TransactionRepo
}

func (u *transactionUsecase) TransferMoney(sender string, receiver string, amount float64) error {
	return u.transactionRepo.TransferMoney(sender, receiver, amount)
}

func (u *transactionUsecase) TopUpBalance(sender string, receiver string, amount float64) error {
	return u.transactionRepo.TopUpBalance(sender, receiver, amount)
}

func (u *transactionUsecase) WithdrawBalance(sender string, receiver string, amount float64) error {
	return u.transactionRepo.WithdrawBalance(sender, receiver, amount)
}

func (u *transactionUsecase) TransferBalance(sender string, receiver string, amount float64) error {
	return u.transactionRepo.TransferBalance(sender, receiver, amount)
}

func NewTransactionUsecase(transactionRepo repository.TransactionRepo) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
	}
}

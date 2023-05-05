package usecase

import (
	"errors"
	"final_project_easycash/repository"
	"final_project_easycash/utils"
	"strconv"
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
	envFilePath := "../.env"
	minTransaction, err := strconv.ParseFloat(utils.DotEnv("MINIMUM_TRANSACTION", envFilePath), 64)
	if err != nil {
		return err
	}
	if amount < minTransaction {
		return errors.New("Minimum Transaction Rp 10.000,00")
	}
	return u.transactionRepo.TransferMoney(sender, receiver, amount)
}

func (u *transactionUsecase) TopUpBalance(sender string, receiver string, amount float64) error {
	envFilePath := "../.env"
	minTransaction, err := strconv.ParseFloat(utils.DotEnv("MINIMUM_TRANSACTION", envFilePath), 64)
	if err != nil {
		return err
	}
	if amount < minTransaction {
		return errors.New("Minimum Transaction Rp 10.000,00")
	}
	adminFee, err := strconv.ParseFloat(utils.DotEnv("ADMIN_FEE_TOPUP", envFilePath), 64)
	if err != nil {
		return err
	}
	amount = amount - adminFee
	return u.transactionRepo.TopUpBalance(sender, receiver, amount)
}

func (u *transactionUsecase) WithdrawBalance(sender string, receiver string, amount float64) error {
	envFilePath := "../.env"
	minTransaction, err := strconv.ParseFloat(utils.DotEnv("MINIMUM_TRANSACTION", envFilePath), 64)
	if err != nil {
		return err
	}
	if amount < minTransaction {
		return errors.New("Minimum Transaction Rp 10.000,00")
	}
	adminFee, err := strconv.ParseFloat(utils.DotEnv("ADMIN_FEE_WITHDRAWAL", envFilePath), 64)
	if err != nil {
		return err
	}
	amount = amount + adminFee
	return u.transactionRepo.WithdrawBalance(sender, receiver, amount)
}

func (u *transactionUsecase) TransferBalance(sender string, receiver string, amount float64) error {
	envFilePath := "../.env"
	minTransaction, err := strconv.ParseFloat(utils.DotEnv("MINIMUM_TRANSACTION", envFilePath), 64)
	if err != nil {
		return err
	}
	if amount < minTransaction {
		return errors.New("Minimum Transaction Rp 10.000,00")
	}
	return u.transactionRepo.TransferBalance(sender, receiver, amount)
}

func NewTransactionUsecase(transactionRepo repository.TransactionRepo) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
	}
}

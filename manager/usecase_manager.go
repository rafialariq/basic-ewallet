package manager

import (
	"final_project_easycash/usecase"
)

type UsecaseManager interface {
	UserUsecase() usecase.UserUsecase
	RegisterUsecase() usecase.RegisterService
	LoginUsecase() usecase.LoginService
	TransactionUsecase() usecase.TransactionUsecase
	HistoryUsecase() usecase.HistoryUsecase
}

type usecaseManager struct {
	repoManager RepoManager
}

func (u *usecaseManager) UserUsecase() usecase.UserUsecase {
	return usecase.NewUserUsecase(u.repoManager.UserRepo(), u.repoManager.FileRepo())
}

func (u *usecaseManager) TransactionUsecase() usecase.TransactionUsecase {
	return usecase.NewTransactionUsecase(u.repoManager.TransactionRepo())
}

func (u *usecaseManager) RegisterUsecase() usecase.RegisterService {
	return usecase.NewRegisterService(u.repoManager.RegisterRepo())
}

func (u *usecaseManager) LoginUsecase() usecase.LoginService {
	return usecase.NewLoginService(u.repoManager.LoginRepo())
}

func (u *usecaseManager) HistoryUsecase() usecase.HistoryUsecase {
	return usecase.NewHistoryUsecase(u.repoManager.HistoryRepo())
}

func NewUsecaseManager(r RepoManager) UsecaseManager {
	return &usecaseManager{
		repoManager: r,
	}
}

package manager

import (
	"final_project_easycash/usecase"
)

type UsecaseManager interface {
	UserUsecase() usecase.UserUsecase
	TransferUsecase() usecase.TransferUsecase
	TopUpUsecase() usecase.TopUpUsecase
	RegisterUsecase() usecase.RegisterService
	LoginUsecase() usecase.LoginService
}

type usecaseManager struct {
	repoManager RepoManager
}

func (u *usecaseManager) UserUsecase() usecase.UserUsecase {
	return usecase.NewUserUsecase(u.repoManager.UserRepo(), u.repoManager.FileRepo())
}

func (u *usecaseManager) TransferUsecase() usecase.TransferUsecase {
	return usecase.NewTransferUsecase(u.repoManager.TransferRepo())
}

func (u *usecaseManager) TopUpUsecase() usecase.TopUpUsecase {
	return usecase.NewTopUpUsecase(u.repoManager.TopUpRepo())
}

func (u *usecaseManager) RegisterUsecase() usecase.RegisterService {
	return usecase.NewRegisterService(u.repoManager.RegisterRepo())
}

func (u *usecaseManager) LoginUsecase() usecase.LoginService {
	return usecase.NewLoginService(u.repoManager.LoginRepo())
}

func NewUsecaseManager(r RepoManager) UsecaseManager {
	return &usecaseManager{
		repoManager: r,
	}
}

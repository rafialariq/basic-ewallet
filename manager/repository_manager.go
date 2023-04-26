package manager

import (
	"final_project_easycash/repository"
)

type RepoManager interface {
	FileRepo() repository.FileRepository
	UserRepo() repository.UserRepo
	TransferRepo() repository.TransferRepo
	TopUpRepo() repository.TopUpRepo
	RegisterRepo() repository.RegisterRepo
	LoginRepo() repository.LoginRepo
}

type repoManager struct {
	infraManager InfraManager
}

func (r *repoManager) FileRepo() repository.FileRepository {
	return repository.NewFileRepository(r.infraManager.InitializeBasePath())
}

func (r *repoManager) UserRepo() repository.UserRepo {
	return repository.NewUserRepo(r.infraManager.ConnectDb())
}

func (r *repoManager) TransferRepo() repository.TransferRepo {
	return repository.NewTransferRepo(r.infraManager.ConnectDb())
}

func (r *repoManager) TopUpRepo() repository.TopUpRepo {
	return repository.NewTopUpRepo(r.infraManager.ConnectDb())
}

func (r *repoManager) RegisterRepo() repository.RegisterRepo {
	return repository.NewRegisterRepo(r.infraManager.ConnectDb())
}

func (r *repoManager) LoginRepo() repository.LoginRepo {
	return repository.NewLoginRepo(r.infraManager.ConnectDb())
}

func NewRepoManager(manager InfraManager) RepoManager {
	return &repoManager{
		infraManager: manager,
	}
}

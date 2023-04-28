package manager

import (
	"final_project_easycash/repository"
)

type RepoManager interface {
	FileRepo() repository.FileRepository
	UserRepo() repository.UserRepo
	RegisterRepo() repository.RegisterRepo
	LoginRepo() repository.LoginRepo
	TransactionRepo() repository.TransactionRepo
	HistoryRepo() repository.HistoryRepo
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

func (r *repoManager) TransactionRepo() repository.TransactionRepo {
	return repository.NewTransactionRepo(r.infraManager.ConnectDb())
}

func (r *repoManager) RegisterRepo() repository.RegisterRepo {
	return repository.NewRegisterRepo(r.infraManager.ConnectDb())
}

func (r *repoManager) LoginRepo() repository.LoginRepo {
	return repository.NewLoginRepo(r.infraManager.ConnectDb())
}

func (r *repoManager) HistoryRepo() repository.HistoryRepo {
	return repository.NewHistoryRepo(r.infraManager.ConnectDb())
}

func NewRepoManager(manager InfraManager) RepoManager {
	return &repoManager{
		infraManager: manager,
	}
}

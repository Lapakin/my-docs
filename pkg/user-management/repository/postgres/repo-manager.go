package postgres

import (
	"github.com/jmoiron/sqlx"

	"github.com/lapotkin/file-storage/pkg/user-management/repository"
)

type RepoManager struct{}

func NewRepoManager() *RepoManager {
	return &RepoManager{}
}

func (rm *RepoManager) NewAuthRepo(db sqlx.ExtContext) repository.AuthRepository {
	return NewAuthRepository(db)
}

func (rm *RepoManager) NewUserRepo(db sqlx.ExtContext) repository.UserRepository {
	return NewUserRepository(db)
}

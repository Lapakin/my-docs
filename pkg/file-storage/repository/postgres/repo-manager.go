package postgres

import (
	"github.com/jmoiron/sqlx"

	"github.com/lapotkin/file-storage/pkg/file-storage/repository"
)

type RepoManager struct{}

func NewRepoManager() *RepoManager {
	return &RepoManager{}
}

func (rm *RepoManager) NewDocumentRepo(db sqlx.ExtContext) repository.DocumentRepository {
	return NewDocumentRepository(db)
}

func (rm *RepoManager) NewFolderRepo(db sqlx.ExtContext) repository.FolderRepository {
	return NewFolderRepository(db)
}

func (rm *RepoManager) NewShareRepo(db sqlx.ExtContext) repository.ShareRepository {
	return NewShareRepository(db)
}

package repository

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"

	"github.com/lapotkin/file-storage/pkg/models"

	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type RepoManager interface {
	NewDocumentRepo(db sqlx.ExtContext) DocumentRepository
	NewFolderRepo(db sqlx.ExtContext) FolderRepository
	NewShareRepo(db sqlx.ExtContext) ShareRepository
}

type DocumentRepository interface {
	Create(ctx context.Context, doc *models.Document) error
	GetByID(ctx context.Context, id uint64) (*models.Document, error)
	FetchDocuments(ctx context.Context, filters f.Filters) (models.Documents, error)
	Update(ctx context.Context, doc *models.Document) error
	Delete(ctx context.Context, id uint64) error
}

type FolderRepository interface {
	Create(ctx context.Context, folder *models.Folder) error
	GetByID(ctx context.Context, id uint64) (*models.Folder, error)
	FetchFolders(ctx context.Context, filters f.Filters) (models.Folders, error)
	Update(ctx context.Context, folder *models.Folder) error
	Delete(ctx context.Context, id uint64) error
	GetPath(ctx context.Context, id uint64) (string, error)
}

type ObjectStorageRepository interface {
	Upload(ctx context.Context, path string, content io.Reader, size int64, contentType string) error
	Download(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
	GetPresignedURL(ctx context.Context, path string) (string, error)
	Copy(ctx context.Context, sourcePath, destPath string) error
}

type ShareRepository interface {
	Create(ctx context.Context, share *models.Share) error
	GetByID(ctx context.Context, id uint64) (*models.Share, error)
	GetByLink(ctx context.Context, link string) (*models.Share, error)
	FetchShares(ctx context.Context, filters f.Filters) (models.Shares, error)
	Update(ctx context.Context, share *models.Share) error
	Delete(ctx context.Context, id uint64) error
	IncrementAccessCount(ctx context.Context, id uint64) error
}

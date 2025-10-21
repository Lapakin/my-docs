package services

import (
	"context"

	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/pkg/models"
)

type Services struct {
	DocumentSvc DocumentSvc
	FolderSvc   FolderSvc
	ShareSvc    ShareSvc
}

type DocumentSvc interface {
	CreateDocument(ctx context.Context, userID uint64, doc *models.Document, content []byte) error
	UpdateDocument(ctx context.Context, doc *models.Document, userID uint64) error
	DeleteDocument(ctx context.Context, id, userID uint64) error
	DownloadDocument(ctx context.Context, id, userID uint64) ([]byte, *models.Document, error)
	GetDocument(ctx context.Context, documentID uint64) (*models.Document, error)
	FetchDocuments(ctx context.Context, userID uint64, folderID *uint64) (models.Documents, error)
	GetPresignedURL(ctx context.Context, id, userID uint64) (string, error)
}

type FolderSvc interface {
	CreateFolder(ctx context.Context, folder *models.Folder) error
	GetFolder(ctx context.Context, id, userID uint64) (*models.Folder, error)
	ListFolders(ctx context.Context, userID uint64, parentID *uint64) (models.Folders, error)
	UpdateFolder(ctx context.Context, folder *models.Folder, userID uint64) error
	DeleteFolder(ctx context.Context, id, userID uint64) error
}

type ShareSvc interface {
	CreateShare(ctx context.Context, claims *jwt.Claims, share *models.Share) error
	GetShareByLink(ctx context.Context, link string) (*models.Share, error)
	DeleteShare(ctx context.Context, id, userID uint64) error
	FetchShares(ctx context.Context, userID uint64) (models.Shares, error)
}

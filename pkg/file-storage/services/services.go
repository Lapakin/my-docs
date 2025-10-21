package services

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/file-storage/repository/minio"
	"github.com/lapotkin/file-storage/pkg/file-storage/repository/postgres"

	pg "github.com/lapotkin/file-storage/internal/adapter/db/postgres"
)

func NewService(db *pg.DB, s3 *s3.Client, rm *postgres.RepoManager, sr *minio.RepoManager, l *logging.Logger) *Services {
	return &Services{
		DocumentSvc: NewDocumentService(db, s3, rm, sr, l.WithField("service", "DocumentService")),
		FolderSvc:   NewFolderService(db, rm, l.WithField("service", "FolderService")),
		ShareSvc:    NewShareService(db, rm, l.WithField("service", "ShareService")),
	}
}

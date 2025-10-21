package services

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lapotkin/file-storage/pkg/utils"

	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/file-storage/repository/minio"
	"github.com/lapotkin/file-storage/pkg/file-storage/repository/postgres"
	"github.com/lapotkin/file-storage/pkg/models"

	pg "github.com/lapotkin/file-storage/internal/adapter/db/postgres"
	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type DocumentService struct {
	db *pg.DB
	s3 *s3.Client
	rm *postgres.RepoManager
	sr *minio.RepoManager
	l  *logging.Logger
}

func NewDocumentService(db *pg.DB, s3 *s3.Client, rm *postgres.RepoManager, sr *minio.RepoManager, l *logging.Logger) *DocumentService {
	return &DocumentService{
		db: db,
		s3: s3,
		rm: rm,
		sr: sr,
		l:  l,
	}
}

func (s *DocumentService) CreateDocument(ctx context.Context, userID uint64, doc *models.Document, content []byte) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error creating document: %v", err)
			s.l.Debugf("document: %+v", doc)
		}
	}()

	now := time.Now()

	doc.UserID = userID
	doc.CreatedAt = now

	s.l.Info("Starting document creation")

	reader := bytes.NewReader(content)

	sr := s.sr.NewObjectStorageRepo(s.s3)
	if err = sr.Upload(ctx, doc.FilePath, reader, doc.FileSize, doc.MimeType); err != nil {
		return ErrInternal
	}

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		if err = sr.Delete(ctx, doc.FilePath); err != nil {
			s.l.Errorf("error cleaning up document after failed DB transaction: %v", err)
		}
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewDocumentRepo(tx).Create(ctx, doc); err != nil {
		if err = sr.Delete(ctx, doc.FilePath); err != nil {
			s.l.Errorf("error cleaning up document after failed DB operation: %v", err)
		}
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		if err = sr.Delete(ctx, doc.FilePath); err != nil {
			s.l.Errorf("error cleaning up document after failed DB commit: %v", err)
		}
		return ErrInternal
	}

	s.l.Info("Successfully created document")
	return nil
}

func (s *DocumentService) GetDocument(ctx context.Context, documentID uint64) (*models.Document, error) {
	doc, err := s.rm.NewDocumentRepo(s.db).GetByID(ctx, documentID)
	if err != nil {
		s.l.Errorf("error fetching document by ID: %v", err)
		s.l.Debugf("documentID: %d", documentID)
		return nil, handleDBError(err)
	}
	return doc, nil
}

func (s *DocumentService) FetchDocuments(ctx context.Context, userID uint64, folderID *uint64) (models.Documents, error) {
	filters := make(f.Filters)
	filters[models.UserIDParam] = utils.UInt64ToString(userID)
	if folderID != nil {
		filters[models.FolderIDParam] = utils.UInt64ToString(*folderID)
	}
	docs, err := s.rm.NewDocumentRepo(s.db).FetchDocuments(ctx, filters)
	if err != nil {
		s.l.Errorf("error fetching documents: %v", err)
		s.l.Debugf("filters: %+v", filters)
		return nil, ErrInternal
	}
	return docs, nil
}

func (s *DocumentService) UpdateDocument(ctx context.Context, doc *models.Document, userID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error updating document: %v", err)
			s.l.Debugf("document: %+v, userID: %d", doc, userID)
		}
	}()

	s.l.Info("Starting document update")

	now := time.Now()

	doc.ModifiedAt = &now

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	r := s.rm.NewDocumentRepo(tx)
	existing, err := r.GetByID(ctx, doc.ID)
	if err != nil {
		return handleDBError(err)
	}

	if existing.UserID != userID {
		return ErrAccessDenied
	}

	if err = r.Update(ctx, doc); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	s.l.Info("Successfully updated document")
	return nil
}

func (s *DocumentService) DeleteDocument(ctx context.Context, id, userID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error deleting document: %v", err)
			s.l.Debugf("documentID: %d, userID: %d", id, userID)
		}
	}()

	s.l.Info("Starting document deletion")

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	doc, err := s.rm.NewDocumentRepo(tx).GetByID(ctx, id)
	if err != nil {
		return handleDBError(err)
	}

	if doc.UserID != userID {
		return ErrAccessDenied
	}

	if err = s.rm.NewDocumentRepo(tx).Delete(ctx, id); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	if err = s.sr.NewObjectStorageRepo(s.s3).Delete(ctx, doc.FilePath); err != nil {
		return ErrInternal
	}

	s.l.Info("Successfully deleted document")
	return nil
}

func (s *DocumentService) DownloadDocument(ctx context.Context, id, userID uint64) ([]byte, *models.Document, error) {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error downloading document: %v", err)
			s.l.Debugf("documentID: %d, userID: %d", id, userID)
		}
	}()

	s.l.Info("Starting document download")

	doc, err := s.rm.NewDocumentRepo(s.db).GetByID(ctx, id)
	if err != nil {
		return nil, nil, handleDBError(err)
	}

	if doc.UserID != userID && !doc.IsPublic {
		return nil, nil, ErrAccessDenied
	}

	reader, err := s.sr.NewObjectStorageRepo(s.s3).Download(ctx, doc.FilePath)
	if err != nil {
		return nil, nil, ErrInternal
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, ErrInternal
	}

	s.l.Info("Successfully downloaded document")
	return content, doc, nil
}

func (s *DocumentService) GetPresignedURL(ctx context.Context, id, userID uint64) (string, error) {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error generating presigned URL: %v", err)
			s.l.Debugf("documentID: %d, userID: %d", id, userID)
		}
	}()

	s.l.Info("Starting presigned URL generation")

	doc, err := s.rm.NewDocumentRepo(s.db).GetByID(ctx, id)
	if err != nil {
		return "", handleDBError(err)
	}

	if doc.UserID != userID && !doc.IsPublic {
		return "", ErrAccessDenied
	}

	url, err := s.sr.NewObjectStorageRepo(s.s3).GetPresignedURL(ctx, doc.FilePath)
	if err != nil {
		return "", ErrInternal
	}

	s.l.Info("Successfully generated presigned URL")
	return url, nil
}

package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/file-storage/repository/postgres"
	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/utils"

	pg "github.com/lapotkin/file-storage/internal/adapter/db/postgres"
	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type ShareService struct {
	db *pg.DB
	rm *postgres.RepoManager
	l  *logging.Logger
}

func NewShareService(db *pg.DB, rm *postgres.RepoManager, l *logging.Logger) *ShareService {
	return &ShareService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *ShareService) CreateShare(ctx context.Context, claims *jwt.Claims, share *models.Share) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error creating share: %v", err)
			s.l.Debugf("share: %+v", share)
		}
	}()

	s.l.Info("Starting share creation")

	share.CreatedAt = time.Now()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	doc, err := s.rm.NewDocumentRepo(tx).GetByID(ctx, share.DocumentID)
	if err != nil {
		return handleDBError(err)
	}

	if doc.UserID != claims.UserID {
		return ErrAccessDenied
	}

	share.OwnerID = claims.UserID

	if share.ShareLink == "" {
		share.ShareLink = uuid.New().String()
	}

	if err = s.rm.NewShareRepo(tx).Create(ctx, share); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	s.l.Info("Successfully created share")
	return nil
}

func (s *ShareService) GetShareByLink(ctx context.Context, link string) (*models.Share, error) {
	share, err := s.rm.NewShareRepo(s.db).GetByLink(ctx, link)
	if err != nil {
		s.l.Errorf("error fetching share by link: %v", err)
		s.l.Debugf("shareLink: %s", link)
		return nil, handleDBError(err)
	}
	return share, nil
}

func (s *ShareService) DeleteShare(ctx context.Context, id, userID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error deleting share: %v", err)
			s.l.Debugf("shareID: %d, userID: %d", id, userID)
		}
	}()

	s.l.Info("Starting share deletion")

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	share, err := s.rm.NewShareRepo(tx).GetByID(ctx, id)
	if err != nil {
		return handleDBError(err)
	}

	if share.OwnerID != userID {
		return ErrAccessDenied
	}

	if err = s.rm.NewShareRepo(tx).Delete(ctx, id); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	s.l.Info("Successfully deleted share")
	return nil
}

func (s *ShareService) FetchShares(ctx context.Context, userID uint64) (models.Shares, error) {
	shares, err := s.rm.NewShareRepo(s.db).FetchShares(ctx, f.Filters{models.OwnerIDParam: utils.UInt64ToString(userID)})
	if err != nil {
		s.l.Errorf("error fetching shares: %v", err)
		s.l.Debugf("userID: %d", userID)
		return nil, ErrInternal
	}
	return shares, nil
}

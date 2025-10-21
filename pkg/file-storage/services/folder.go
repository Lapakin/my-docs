package services

import (
	"context"
	"time"

	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/file-storage/repository/postgres"
	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/utils"

	pg "github.com/lapotkin/file-storage/internal/adapter/db/postgres"
	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type FolderService struct {
	db *pg.DB
	rm *postgres.RepoManager
	l  *logging.Logger
}

func NewFolderService(db *pg.DB, rm *postgres.RepoManager, l *logging.Logger) *FolderService {
	return &FolderService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *FolderService) CreateFolder(ctx context.Context, folder *models.Folder) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error creating folder: %v", err)
			s.l.Debugf("folder: %+v", folder)
		}
	}()

	s.l.Info("Starting folder creation")

	folder.CreatedAt = time.Now()

	r := s.rm.NewFolderRepo(s.db)
	if folder.ParentID != nil {
		var parentPath string
		parentPath, err = r.GetPath(ctx, *folder.ParentID)
		if err != nil {
			return handleDBError(err)
		}
		folder.Path = parentPath + "/" + folder.Name
	} else {
		folder.Path = "/" + folder.Name
	}

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewFolderRepo(tx).Create(ctx, folder); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	s.l.Info("Successfully created folder")
	return nil
}

func (s *FolderService) GetFolder(ctx context.Context, id, userID uint64) (*models.Folder, error) {
	folder, err := s.rm.NewFolderRepo(s.db).GetByID(ctx, id)
	if err != nil {
		s.l.Errorf("error fetching folder by ID: %v", err)
		s.l.Debugf("folderID: %d", id)
		return nil, handleDBError(err)
	}

	if folder.UserID != userID {
		return nil, ErrAccessDenied
	}

	return folder, nil
}

func (s *FolderService) ListFolders(ctx context.Context, userID uint64, parentID *uint64) (models.Folders, error) {
	filters := make(f.Filters)
	filters[models.UserIDParam] = utils.UInt64ToString(userID)
	if parentID != nil {
		filters[models.ParentIDParam] = utils.UInt64ToString(*parentID)
	}
	folders, err := s.rm.NewFolderRepo(s.db).FetchFolders(ctx, filters)
	if err != nil {
		s.l.Errorf("error fetching folders: %v", err)
		s.l.Debugf("filters: %+v", filters)
		return nil, ErrInternal
	}
	return folders, nil
}

func (s *FolderService) UpdateFolder(ctx context.Context, folder *models.Folder, userID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error updating folder: %v", err)
			s.l.Debugf("folder: %+v, userID: %d", folder, userID)
		}
	}()

	s.l.Info("Starting folder update")

	now := time.Now()

	folder.ModifiedAt = &now

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	r := s.rm.NewFolderRepo(tx)
	existing, err := r.GetByID(ctx, folder.ID)
	if err != nil {
		return handleDBError(err)
	}

	if existing.UserID != userID {
		return ErrAccessDenied
	}

	if err = r.Update(ctx, folder); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	s.l.Info("Successfully updated folder")
	return nil
}

func (s *FolderService) DeleteFolder(ctx context.Context, id, userID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error deleting folder: %v", err)
			s.l.Debugf("folderID: %d, userID: %d", id, userID)
		}
	}()

	s.l.Info("Starting folder deletion")

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	r := s.rm.NewFolderRepo(tx)
	folder, err := r.GetByID(ctx, id)
	if err != nil {
		return handleDBError(err)
	}

	if folder.UserID != userID {
		return ErrAccessDenied
	}

	if err = r.Delete(ctx, id); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	s.l.Info("Successfully deleted folder")
	return nil
}

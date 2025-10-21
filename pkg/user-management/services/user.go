package services

import (
	"context"

	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/user-management/repository"

	pg "github.com/lapotkin/file-storage/internal/adapter/db/postgres"
	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type UserService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewUserService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *UserService {
	return &UserService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, userID uint64) (*models.User, error) {
	user, err := s.rm.NewUserRepo(s.db).GetUserByID(ctx, userID)
	if err != nil {
		s.l.Errorf("error fetching user by ID: %v", err)
		s.l.Debugf("userID: %d", userID)
		return nil, handleDBError(err)
	}
	return user, nil
}

func (s *UserService) FetchUsers(ctx context.Context, filters f.Filters) (models.Users, error) {
	users, err := s.rm.NewUserRepo(s.db).FetchUsers(ctx, filters)
	if err != nil {
		s.l.Errorf("error fetching users: %v", err)
		s.l.Debugf("filters: %+v", filters)
		return nil, ErrInternal
	}
	return users, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error deleting user: %v", err)
			s.l.Debugf("userID: %d", userID)
		}
	}()

	s.l.Debug("Starting user deletion")

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewUserRepo(tx).DeleteUser(ctx, userID); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	s.l.Debug("Successfully deleted user")
	return nil
}

func (s *UserService) DeactivateUser(ctx context.Context, userID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error deactivating user: %v", err)
			s.l.Debugf("userID: %d", userID)
		}
	}()

	s.l.Debug("Starting user deactivation")

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewUserRepo(tx).DeactivateUser(ctx, userID); err != nil {
		return ErrInternal
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	s.l.Debug("Successfully deactivated user")
	return nil
}

func (s *UserService) ActivateUser(ctx context.Context, userID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error activating user: %v", err)
			s.l.Debugf("userID: %d", userID)
		}
	}()

	s.l.Debug("Starting user activation")

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewUserRepo(tx).ActivateUser(ctx, userID); err != nil {
		return ErrInternal
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	s.l.Debug("Successfully activated user")
	return nil
}

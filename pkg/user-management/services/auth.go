package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/user-management/repository"

	pg "github.com/lapotkin/file-storage/internal/adapter/db/postgres"
)

type AuthService struct {
	db *pg.DB
	rm repository.RepoManager
	tm *jwt.TokenManager
	l  *logging.Logger
}

func NewAuthService(db *pg.DB, rm repository.RepoManager, tm *jwt.TokenManager, l *logging.Logger) *AuthService {
	return &AuthService{
		db: db,
		rm: rm,
		tm: tm,
		l:  l,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during registration: %v", err)
			s.l.Debugf("req: %+v", req)
		}
	}()

	s.l.Info("Starting user registration")

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return nil, ErrInternal
	}
	defer tx.Rollback()

	r := s.rm.NewUserRepo(tx)
	dbUser, err := r.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, handleDBError(err)
	}

	if dbUser != nil {
		return nil, ErrUserExists
	}

	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.RoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err = r.CreateUser(ctx, user); err != nil {
		return nil, handleDBError(err)
	}

	password, err := user.SetPassword(req.Password)
	if err != nil {
		return nil, ErrInternal
	}

	passwordData := &models.PasswordData{
		UserID:     user.ID,
		Hash:       password.PasswordHash,
		ModifiedAt: time.Now(),
	}

	authRepo := s.rm.NewAuthRepo(tx)
	if err = authRepo.CreatePassword(ctx, passwordData); err != nil {
		return nil, handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, ErrInternal
	}

	s.l.Info("Successfully registered new user")
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during login: %v", err)
			s.l.Debugf("req: %+v", req)
		}
	}()

	s.l.Info("Starting user login")

	userRepo := s.rm.NewUserRepo(s.db)
	user, err := userRepo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, handleDBError(err)
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserDeactivated
	}

	authRepo := s.rm.NewAuthRepo(s.db)
	passwordHash, err := authRepo.GetPasswordHash(ctx, user.ID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !models.CheckPassword(passwordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tm.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Username,
		string(user.Role),
	)
	if err != nil {
		return nil, ErrInternal
	}

	s.l.Info("User logged in successfully")
	return &models.LoginResponse{AccessToken: token}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uint64, req *models.ChangePasswordRequest) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during password change: %v", err)
			s.l.Debugf("userID: %d", userID)
		}
	}()

	s.l.Info("Starting password change")

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	r := s.rm.NewAuthRepo(tx)
	currentHash, err := r.GetPasswordHash(ctx, userID)
	if err != nil {
		return handleDBError(err)
	}

	if !models.CheckPassword(currentHash, req.OldPassword) {
		return ErrInvalidCredentials
	}

	newHash, err := models.HashPassword(req.NewPassword)
	if err != nil {
		return ErrInternal
	}

	password := &models.PasswordData{
		UserID:     userID,
		Hash:       newHash,
		ModifiedAt: time.Now(),
	}

	if err = r.UpdatePassword(ctx, password); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	s.l.Info("Password changed successfully")
	return nil
}

package services

import (
	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/user-management/repository"

	pg "github.com/lapotkin/file-storage/internal/adapter/db/postgres"
)

func NewService(db *pg.DB, rm repository.RepoManager, tm *jwt.TokenManager, l *logging.Logger) *Services {
	return &Services{
		AuthSvc: NewAuthService(db, rm, tm, l.WithField("service", "AuthService")),
		UserSvc: NewUserService(db, rm, l.WithField("service", "UserService")),
	}
}

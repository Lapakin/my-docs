package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lapotkin/file-storage/pkg/models"

	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type RepoManager interface {
	NewAuthRepo(db sqlx.ExtContext) AuthRepository
	NewUserRepo(db sqlx.ExtContext) UserRepository
}

type AuthRepository interface {
	GetRefreshToken(ctx context.Context, token string, currentTime time.Time) (uint64, error)
	SaveRefreshToken(ctx context.Context, data *models.RefreshTokenData) error
	RevokeRefreshToken(ctx context.Context, token string) error
	GetPasswordHash(ctx context.Context, userID uint64) (string, error)
	CreatePassword(ctx context.Context, data *models.PasswordData) error
	UpdatePassword(ctx context.Context, data *models.PasswordData) error
}
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID uint64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	FetchUsers(ctx context.Context, filters f.Filters) (models.Users, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID uint64) error
	HardDeleteUser(ctx context.Context, userID uint64) error
	DeactivateUser(ctx context.Context, userID uint64) error
	ActivateUser(ctx context.Context, userID uint64) error
}

package services

import (
	"context"

	"github.com/lapotkin/file-storage/pkg/models"

	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type Services struct {
	AuthSvc AuthSvc
	UserSvc UserSvc
}

type AuthSvc interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	ChangePassword(ctx context.Context, userID uint64, req *models.ChangePasswordRequest) error
}

type UserSvc interface {
	GetUserByID(ctx context.Context, userID uint64) (*models.User, error)
	FetchUsers(ctx context.Context, filters f.Filters) (models.Users, error)
	DeleteUser(ctx context.Context, userID uint64) error
	DeactivateUser(ctx context.Context, userID uint64) error
	ActivateUser(ctx context.Context, userID uint64) error
}

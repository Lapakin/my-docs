package services

import (
	"github.com/lapotkin/file-storage/internal/adapter/krakend"
	"github.com/lapotkin/file-storage/pkg/models"
)

type Services struct {
	AuthSvc       AuthService
	UserSvc       UserService
	DocumentSvc   DocumentService
	FolderSvc     FolderService
	KrakendClient *krakend.Client
}

type AuthService interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	ValidateToken(token string) (*models.User, error)
}

type UserService interface {
	GetUser(userID string, token string) (*models.User, error)
	UpdateUser(userID string, user *models.User, token string) (*models.User, error)
	DeleteUser(userID string, token string) error
	ListUsers(search, limit, offset, token string) (models.Users, error)
}

type FolderService interface {
	CreateFolder(name, parentID string, token string) (*models.Folder, error)
	GetFolder(folderID string, token string) (*models.Folder, error)
	ListFolders(parentID string, token string) (models.Folders, error)
	UpdateFolder(folderID string, folder *models.Folder, token string) error
	DeleteFolder(folderID string, token string) error
}

type DocumentService interface {
	CreateDocument(name, folderID string, token string) (*models.Document, error)
	GetDocument(documentID string, token string) (*models.Document, error)
	ListDocuments(folderID string, token string) (models.Documents, error)
	UpdateDocument(documentID string, doc *models.Document, token string) error
	DeleteDocument(documentID string, token string) error
	SearchDocuments(query string, token string) (models.Documents, error)
}

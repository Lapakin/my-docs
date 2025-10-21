package services

import (
	"github.com/lapotkin/file-storage/internal/adapter/krakend"
)

func NewServices(krakendClient *krakend.Client) *Services {
	return &Services{
		AuthSvc:       NewAuthService(krakendClient),
		UserSvc:       NewUserService(krakendClient),
		DocumentSvc:   NewDocumentService(krakendClient),
		FolderSvc:     NewFolderService(krakendClient),
		KrakendClient: krakendClient,
	}
}

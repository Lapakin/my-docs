package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lapotkin/file-storage/internal/adapter/http/middleware"
	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/file-storage/handlers"
	"github.com/lapotkin/file-storage/pkg/file-storage/services"
)

func NewRouter(svc *services.Services, tm *jwt.TokenManager, l *logging.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggerWithCustomLogger(l))
	r.Use(gin.Recovery())

	r.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.Auth(tm))

	apiV1.Handle(http.MethodPost, "/documents", handlers.NewUploadDocumentHandler(svc.DocumentSvc))
	apiV1.Handle(http.MethodPut, "/documents/:id", handlers.NewUpdateDocumentHandler(svc.DocumentSvc))
	apiV1.Handle(http.MethodDelete, "/documents/:id", handlers.NewDeleteDocumentHandler(svc.DocumentSvc))
	apiV1.Handle(http.MethodGet, "/documents/:id", handlers.NewGetDocumentHandler(svc.DocumentSvc))
	apiV1.Handle(http.MethodGet, "/documents/:id/download", handlers.NewDownloadDocumentHandler(svc.DocumentSvc))
	apiV1.Handle(http.MethodGet, "/documents", handlers.NewFetchDocumentsHandler(svc.DocumentSvc))

	apiV1.Handle(http.MethodPost, "/folders", handlers.NewCreateFolderHandler(svc.FolderSvc))
	apiV1.Handle(http.MethodPut, "/folders/:id", handlers.NewUpdateFolderHandler(svc.FolderSvc))
	apiV1.Handle(http.MethodDelete, "/folders/:id", handlers.NewDeleteFolderHandler(svc.FolderSvc))
	apiV1.Handle(http.MethodGet, "/folders/:id", handlers.NewGetFolderHandler(svc.FolderSvc))
	apiV1.Handle(http.MethodGet, "/folders", handlers.NewListFoldersHandler(svc.FolderSvc))

	apiV1.Handle(http.MethodPost, "/shares", handlers.NewCreateShareHandler(svc.ShareSvc))
	apiV1.Handle(http.MethodGet, "/shares/:id", handlers.NewGetShareByLinkHandler(svc.ShareSvc))
	apiV1.Handle(http.MethodDelete, "/shares/:id", handlers.NewDeleteShareHandler(svc.ShareSvc))
	apiV1.Handle(http.MethodGet, "/shares", handlers.NewFetchSharesHandler(svc.ShareSvc))

	return r
}

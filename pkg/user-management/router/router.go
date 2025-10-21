package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lapotkin/file-storage/internal/adapter/http/middleware"
	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/user-management/handlers"
	"github.com/lapotkin/file-storage/pkg/user-management/services"
)

func NewRouter(svc *services.Services, tm *jwt.TokenManager, logger *logging.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggerWithCustomLogger(logger))
	r.Use(gin.Recovery())

	r.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	apiV1 := r.Group("/api/v1")
	apiV1.Handle(http.MethodPost, "/auth/register", handlers.NewRegisterHandler(svc.AuthSvc))
	apiV1.Handle(http.MethodPost, "/auth/login", handlers.NewLoginHandler(svc.AuthSvc))

	protected := apiV1.Group("")
	protected.Use(middleware.Auth(tm))
	{
		protected.Handle(http.MethodPost, "/auth/change-password", handlers.NewChangePasswordHandler(svc.AuthSvc))
		protected.Handle(http.MethodGet, "/users/:id", handlers.NewGetUserByIDHandler(svc.UserSvc))
		protected.Handle(http.MethodGet, "/users", handlers.NewFetchUsersHandler(svc.UserSvc))
	}

	admin := apiV1.Group("")
	admin.Use(middleware.Auth(tm))
	admin.Use(middleware.Admin())
	{
		admin.Handle(http.MethodDelete, "/users/:id", handlers.NewDeleteUserHandler(svc.UserSvc))
		admin.Handle(http.MethodPost, "/users/:id/deactivate", handlers.NewDeactivateUserHandler(svc.UserSvc))
		admin.Handle(http.MethodPost, "/users/:id/activate", handlers.NewActivateUserHandler(svc.UserSvc))
	}

	return r
}

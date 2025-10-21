package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/lapotkin/file-storage/internal/adapter/http/middleware"
	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/ui/handlers"
	"github.com/lapotkin/file-storage/pkg/ui/services"
)

func NewRouter(svc *services.Services, logger *logging.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggerWithCustomLogger(logger))
	r.Use(gin.Recovery())

	r.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.HEAD("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	templatePath := "templates/*.html"
	if _, err := os.Stat("pkg/ui/templates"); err == nil {
		templatePath = "pkg/ui/templates/*.html"
	}
	r.LoadHTMLGlob(templatePath)

	staticPath := "static"
	if _, err := os.Stat("pkg/ui/static"); err == nil {
		staticPath = "pkg/ui/static"
	}
	r.Static("/static", staticPath)

	r.GET("/", handlers.NewIndexPageHandler())
	r.GET("/login", handlers.NewLoginPageHandler())
	r.GET("/register", handlers.NewRegisterPageHandler())

	r.POST("/register", handlers.NewRegisterHandler(svc.AuthSvc))
	r.POST("/login", handlers.NewLoginHandler(svc.AuthSvc))

	protected := r.Group("/")
	protected.Use(middleware.CookieAuthMiddleware())
	{
		protected.GET("/dashboard", handlers.NewDashboardPageHandler())
		protected.GET("/documents", handlers.NewDocumentsPageHandler())
		protected.GET("/folders", handlers.NewFoldersPageHandler())
		protected.GET("/shares", handlers.NewSharesPageHandler())
		protected.GET("/profile", handlers.NewProfilePageHandler())
	}

	apiV1 := r.Group("/api")
	{
		apiV1.POST("/login", handlers.NewLoginHandler(svc.AuthSvc))
		apiV1.POST("/register", handlers.NewRegisterHandler(svc.AuthSvc))
		apiV1.GET("/logout", handlers.NewLogoutHandler())

		apiProtected := apiV1.Group("")
		apiProtected.Use(middleware.CookieAuthMiddleware())
		{
			apiProtected.GET("/users", handlers.NewListUsersHandler(svc.UserSvc))
			apiProtected.GET("/users/:id", handlers.NewUserProfileHandler(svc.UserSvc))
			apiProtected.PUT("/users/:id", handlers.NewUpdateUserHandler(svc.UserSvc))

			apiProtected.POST("/documents/upload", handlers.NewDocumentUploadProxyHandler(svc.KrakendClient))
			apiProtected.GET("/documents/:id/download", handlers.NewDocumentDownloadProxyHandler(svc.KrakendClient))

			apiProtected.Any("/v1/*path", handlers.NewAPIProxyHandler(svc.KrakendClient, "/api/v1"))
		}
	}

	return r
}

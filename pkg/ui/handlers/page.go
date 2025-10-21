package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewIndexPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home",
		})
	}
}

func NewLoginPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := gin.H{
			"title": "Login",
		}

		if c.Query("registered") == "true" {
			data["success"] = "Registration successful! Please log in."
		}

		c.HTML(http.StatusOK, "login.html", data)
	}
}

func NewRegisterPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Register",
		})
	}
}

func NewDashboardPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "Dashboard",
		})
	}
}

func NewDocumentsPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "documents.html", gin.H{
			"title": "Documents",
		})
	}
}

func NewFoldersPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "folders.html", gin.H{
			"title": "Folders",
		})
	}
}

func NewSharesPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "shares.html", gin.H{
			"title": "Shares",
		})
	}
}

func NewProfilePageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"title": "Profile",
		})
	}
}

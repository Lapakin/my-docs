package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/ui/services"
)

func NewLoginHandler(svc services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest

		contentType := c.GetHeader("Content-Type")
		isJSON := strings.Contains(contentType, "application/json")

		if isJSON {
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
		} else {
			if err := c.ShouldBind(&req); err != nil {
				c.HTML(http.StatusBadRequest, "login.html", gin.H{
					"title": "Login",
					"error": "Invalid form data: " + err.Error(),
				})
				return
			}
		}

		loginResp, err := svc.Login(&req)
		if err != nil {
			if errors.Is(err, services.ErrUnauthorized) {
				if isJSON {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
					return
				}
				c.HTML(http.StatusUnauthorized, "login.html", gin.H{
					"title": "Login",
					"error": "Invalid email or password",
				})
				return
			}
			if isJSON {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"title": "Login",
				"error": "Login failed. Please try again.",
			})
			return
		}

		c.SetCookie("token", loginResp.AccessToken, 3600*24, "/", "", false, true)

		if isJSON {
			c.JSON(http.StatusOK, gin.H{"access_token": loginResp.AccessToken, "message": "login successful"})
		} else {
			c.Redirect(http.StatusFound, "/dashboard")
		}
	}
}

func NewRegisterHandler(svc services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest

		contentType := c.GetHeader("Content-Type")
		isJSON := strings.Contains(contentType, "application/json")

		if isJSON {
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
				return
			}
		} else {
			if err := c.ShouldBind(&req); err != nil {
				c.HTML(http.StatusBadRequest, "register.html", gin.H{
					"title": "Register",
					"error": "Invalid form data: " + err.Error(),
				})
				return
			}
		}

		user, err := svc.Register(&req)
		if err != nil {
			var apiErr *services.APIError
			if errors.As(err, &apiErr) {
				if isJSON {
					c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
					return
				}
				c.HTML(apiErr.StatusCode, "register.html", gin.H{
					"title": "Register",
					"error": apiErr.Message,
				})
				return
			}
			if isJSON {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}
			c.HTML(http.StatusInternalServerError, "register.html", gin.H{
				"title": "Register",
				"error": "Registration failed. Please try again.",
			})
			return
		}

		if isJSON {
			c.JSON(http.StatusCreated, gin.H{
				"message": "registration successful",
				"user":    user,
			})
		} else {
			c.Redirect(http.StatusFound, "/login?registered=true")
		}
	}
}

func NewLogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "", false, true)

		if c.GetHeader("Accept") == "application/json" {
			c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}

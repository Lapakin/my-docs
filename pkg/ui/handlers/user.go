package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/ui/services"
)

func NewUserProfileHandler(userSvc services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		userID := c.Param(models.IDParam)
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
			return
		}

		user, err := userSvc.GetUser(userID, token)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			if errors.Is(err, services.ErrUnauthorized) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func NewListUsersHandler(userSvc services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		search := c.Query("search")
		limit := c.DefaultQuery("limit", "20")
		offset := c.DefaultQuery("offset", "0")

		users, err := userSvc.ListUsers(search, limit, offset, token)
		if err != nil {
			if errors.Is(err, services.ErrUnauthorized) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func NewUpdateUserHandler(userSvc services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		userID := c.Param(models.IDParam)
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
			return
		}

		var updateReq models.User
		if err = c.ShouldBindJSON(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		updatedUser, err := userSvc.UpdateUser(userID, &updateReq, token)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			if errors.Is(err, services.ErrUnauthorized) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "user updated successfully",
			"user":    updatedUser,
		})
	}
}

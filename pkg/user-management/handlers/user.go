package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lapotkin/file-storage/internal/adapter/http/rest"
	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/user-management/services"
	"github.com/lapotkin/file-storage/pkg/utils"
)

func NewGetUserByIDHandler(svc services.UserSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		user, err := svc.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, user)
	}
}

func NewFetchUsersHandler(svc services.UserSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters, err := rest.CreateFiltersFromQuery(c.Request.URL.Query(), rest.Queries{
			{Param: models.IDsParam, ValidateFunc: nil},
			{Param: models.EmailParam, ValidateFunc: nil},
		})
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrParseQuery)
			return
		}

		users, err := svc.FetchUsers(c.Request.Context(), filters)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, users)
	}
}

func NewDeleteUserHandler(svc services.UserSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		if err = svc.DeleteUser(c.Request.Context(), userID); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func NewDeactivateUserHandler(svc services.UserSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		if err = svc.DeactivateUser(c.Request.Context(), userID); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func NewActivateUserHandler(svc services.UserSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		if err = svc.ActivateUser(c.Request.Context(), userID); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lapotkin/file-storage/internal/adapter/http/rest"
	"github.com/lapotkin/file-storage/internal/domain/json"
	"github.com/lapotkin/file-storage/pkg/file-storage/services"
	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/utils"
)

func NewCreateFolderHandler(svc services.FolderSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}
		defer c.Request.Body.Close()

		var folder *models.Folder
		if err = json.Unmarshal(body, &folder); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if folder == nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		folder.UserID = userID

		if err = svc.CreateFolder(c.Request.Context(), folder); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusCreated, folder)
	}
}

func NewGetFolderHandler(svc services.FolderSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		folderID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		folder, err := svc.GetFolder(c.Request.Context(), folderID, userID)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, folder)
	}
}

func NewListFoldersHandler(svc services.FolderSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		var parentID *uint64
		if pid := c.Query("parent_id"); pid != "" {
			var id uint64
			id, err = utils.StringToUint64(pid)
			if err != nil {
				rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
				return
			}
			parentID = &id
		}

		folders, err := svc.ListFolders(c.Request.Context(), userID, parentID)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, folders)
	}
}

func NewUpdateFolderHandler(svc services.FolderSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}
		defer c.Request.Body.Close()

		var folder *models.Folder
		if err = json.Unmarshal(body, &folder); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if folder == nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.UpdateFolder(c.Request.Context(), folder, userID); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, folder)
	}
}

func NewDeleteFolderHandler(svc services.FolderSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		folderID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		if err = svc.DeleteFolder(c.Request.Context(), folderID, userID); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lapotkin/file-storage/internal/adapter/http/rest"
	"github.com/lapotkin/file-storage/internal/domain/json"
	"github.com/lapotkin/file-storage/pkg/file-storage/services"
	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/utils"
)

func NewUploadDocumentHandler(svc services.DocumentSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, rest.ErrConvID)
			return
		}

		if err = c.Request.ParseMultipartForm(models.MaxUploadSize); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrFileTooLarge)
			return
		}

		file, header, err := c.Request.FormFile(models.FileParam)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrFileRequired)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}

		name := c.PostForm("name")
		if name == "" {
			name = header.Filename
		}

		description := c.PostForm("description")
		folderIDStr := c.PostForm("folder_id")

		doc := &models.Document{
			Name:        name,
			Description: description,
			FilePath:    filepath.Join("documents", strconv.FormatUint(userID, 10), uuid.New().String()+filepath.Ext(header.Filename)),
			FileSize:    header.Size,
			MimeType:    header.Header.Get("Content-Type"),
			IsPublic:    c.PostForm("is_public") == "true",
		}

		if folderIDStr != "" {
			var folderID uint64
			folderID, err = utils.StringToUint64(folderIDStr)
			if err != nil {
				rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
				return
			}
			doc.FolderID = &folderID
		}

		if err = svc.CreateDocument(c.Request.Context(), userID, doc, content); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusCreated, doc)
	}
}

func NewGetDocumentHandler(svc services.DocumentSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		doc, err := svc.GetDocument(c.Request.Context(), documentID)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, doc)
	}
}

func NewFetchDocumentsHandler(svc services.DocumentSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		var folderID *uint64
		if fid := c.Query("folder_id"); fid != "" {
			id, convErr := utils.StringToUint64(fid)
			if convErr != nil {
				rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
				return
			}
			folderID = &id
		}

		docs, err := svc.FetchDocuments(c.Request.Context(), userID, folderID)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, docs)
	}
}

func NewUpdateDocumentHandler(svc services.DocumentSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		documentID, err := utils.StringToUint64(c.Param(models.IDParam))
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

		var doc *models.Document
		if err = json.Unmarshal(body, &doc); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}
		doc.ID = documentID

		if doc == nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.UpdateDocument(c.Request.Context(), doc, userID); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, doc)
	}
}

func NewDeleteDocumentHandler(svc services.DocumentSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		if err = svc.DeleteDocument(c.Request.Context(), documentID, userID); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func NewDownloadDocumentHandler(svc services.DocumentSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		content, doc, err := svc.DownloadDocument(c.Request.Context(), documentID, userID)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		c.Header("Content-Disposition", "attachment; filename="+doc.Name)
		c.Header("Content-Type", doc.MimeType)
		c.Header("Content-Length", strconv.FormatInt(doc.FileSize, 10))
		c.Data(http.StatusOK, doc.MimeType, content)
	}
}

func NewGetPresignedURLHandler(svc services.DocumentSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		url, err := svc.GetPresignedURL(c.Request.Context(), documentID, userID)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, gin.H{"url": url})
	}
}

package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lapotkin/file-storage/internal/adapter/http/rest"
	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/internal/domain/json"
	"github.com/lapotkin/file-storage/pkg/file-storage/services"
	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/utils"
)

func NewCreateShareHandler(svc services.ShareSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			rest.RespondError(c, http.StatusUnauthorized, rest.ErrUnauthorized)
			return
		}

		jwtClaims, ok := claims.(*jwt.Claims)
		if !ok {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrBadRequest)
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}

		var share *models.Share
		if err = json.Unmarshal(body, &share); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if share == nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.CreateShare(c.Request.Context(), jwtClaims, share); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusCreated, share)
	}
}

func NewGetShareByLinkHandler(svc services.ShareSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		link := c.Param("link")
		if link == "" {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrParamRequired)
			return
		}

		share, err := svc.GetShareByLink(c.Request.Context(), link)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, share)
	}
}

func NewFetchSharesHandler(svc services.ShareSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, rest.ErrConvID)
			return
		}

		shares, err := svc.FetchShares(c.Request.Context(), userID)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, shares)
	}
}

func NewDeleteShareHandler(svc services.ShareSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		shareID, err := utils.StringToUint64(c.Param(models.IDParam))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, rest.ErrConvID)
			return
		}

		if err = svc.DeleteShare(c.Request.Context(), shareID, userID); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

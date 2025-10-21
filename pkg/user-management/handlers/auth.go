package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lapotkin/file-storage/internal/adapter/http/rest"
	"github.com/lapotkin/file-storage/internal/domain/json"
	"github.com/lapotkin/file-storage/pkg/models"
	"github.com/lapotkin/file-storage/pkg/user-management/services"
	"github.com/lapotkin/file-storage/pkg/utils"
)

func NewRegisterHandler(svc services.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}
		defer c.Request.Body.Close()

		var req *models.RegisterRequest
		if err = json.Unmarshal(body, &req); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if req == nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = req.Validate(); err != nil {
			rest.RespondError(c, http.StatusBadRequest, err)
			return
		}

		user, err := svc.Register(c.Request.Context(), req)
		if err != nil {
			if errors.Is(err, services.ErrUserExists) {
				rest.RespondError(c, http.StatusConflict, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusCreated, user)
	}
}

func NewLoginHandler(svc services.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}
		defer c.Request.Body.Close()

		var req *models.LoginRequest
		if err = json.Unmarshal(body, &req); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if req == nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = req.Validate(); err != nil {
			rest.RespondError(c, http.StatusBadRequest, err)
			return
		}

		token, err := svc.Login(c.Request.Context(), req)
		if err != nil {
			if errors.Is(err, services.ErrInvalidCredentials) {
				rest.RespondError(c, http.StatusUnauthorized, err)
				return
			}
			if errors.Is(err, services.ErrUserDeactivated) {
				rest.RespondError(c, http.StatusForbidden, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, token)
	}
}

func NewChangePasswordHandler(svc services.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.GetString(models.UserIDParam))
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, rest.ErrConvID)
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}
		defer c.Request.Body.Close()

		var req *models.ChangePasswordRequest
		if err = json.Unmarshal(body, &req); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if req == nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = req.Validate(); err != nil {
			rest.RespondError(c, http.StatusBadRequest, err)
			return
		}

		if err = svc.ChangePassword(c.Request.Context(), userID, req); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

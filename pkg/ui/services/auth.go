package services

import (
	"fmt"
	"net/http"

	"github.com/lapotkin/file-storage/internal/adapter/krakend"
	"github.com/lapotkin/file-storage/pkg/models"
)

type authService struct {
	kc *krakend.Client
}

func NewAuthService(kc *krakend.Client) AuthService {
	return &authService{
		kc: kc,
	}
}

func (s *authService) Register(req *models.RegisterRequest) (*models.User, error) {
	if err := req.Validate(); err != nil {
		return nil, NewAPIError(http.StatusBadRequest, "invalid request", err)
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodPost,
		Path:   "/api/v1/auth/register",
		Body:   req,
	})
	if err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to send request", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, NewAPIError(resp.StatusCode, "registration failed", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var user models.User
	if err = resp.DecodeJSON(&user); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return &user, nil
}

func (s *authService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, NewAPIError(http.StatusBadRequest, "invalid request", err)
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodPost,
		Path:   "/api/v1/auth/login",
		Body:   req,
	})
	if err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to send request", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, NewAPIError(resp.StatusCode, "login failed", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var loginResp models.LoginResponse
	if err = resp.DecodeJSON(&loginResp); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return &loginResp, nil
}

func (s *authService) ValidateToken(token string) (*models.User, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodGet,
		Path:   "/api/v1/token/validate",
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	if err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to send request", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, NewAPIError(resp.StatusCode, "token validation failed", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var user models.User
	if err = resp.DecodeJSON(&user); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return &user, nil
}

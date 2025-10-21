package services

import (
	"fmt"
	"net/http"

	"github.com/lapotkin/file-storage/internal/adapter/krakend"
	"github.com/lapotkin/file-storage/pkg/models"
)

type userService struct {
	kc *krakend.Client
}

func NewUserService(kc *krakend.Client) UserService {
	return &userService{
		kc: kc,
	}
}

func (s *userService) GetUser(userID string, token string) (*models.User, error) {
	if userID == "" {
		return nil, NewAPIError(http.StatusBadRequest, "user ID is required", nil)
	}
	if token == "" {
		return nil, ErrUnauthorized
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/api/v1/users/%s", userID),
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	if err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to send request", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, NewAPIError(resp.StatusCode, "failed to get user", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var user models.User
	if err = resp.DecodeJSON(&user); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return &user, nil
}

func (s *userService) UpdateUser(userID string, user *models.User, token string) (*models.User, error) {
	if userID == "" {
		return nil, NewAPIError(http.StatusBadRequest, "user ID is required", nil)
	}
	if token == "" {
		return nil, ErrUnauthorized
	}
	if user == nil {
		return nil, NewAPIError(http.StatusBadRequest, "user data is required", nil)
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/api/v1/users/%s", userID),
		Body:   user,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	if err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to send request", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, NewAPIError(resp.StatusCode, "failed to update user", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var updatedUser models.User
	if err = resp.DecodeJSON(&updatedUser); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return &updatedUser, nil
}

func (s *userService) DeleteUser(userID string, token string) error {
	if userID == "" {
		return NewAPIError(http.StatusBadRequest, "user ID is required", nil)
	}
	if token == "" {
		return ErrUnauthorized
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/api/v1/users/%s", userID),
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "failed to send request", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return NewAPIError(resp.StatusCode, "failed to delete user", fmt.Errorf("response: %s", string(resp.Body)))
	}

	return nil
}

func (s *userService) ListUsers(search, limit, offset, token string) (models.Users, error) {
	if token == "" {
		return nil, ErrUnauthorized
	}

	path := "/api/v1/users"
	if search != "" || limit != "" || offset != "" {
		path += "?"
		params := make([]string, 0)
		if search != "" {
			params = append(params, "search="+search)
		}
		if limit != "" {
			params = append(params, "limit="+limit)
		}
		if offset != "" {
			params = append(params, "offset="+offset)
		}
		path += fmt.Sprintf("%s", params[0])
		for i := 1; i < len(params); i++ {
			path += "&" + params[i]
		}
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodGet,
		Path:   path,
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
		return nil, NewAPIError(resp.StatusCode, "failed to list users", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var users models.Users
	if err = resp.DecodeJSON(&users); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return users, nil
}

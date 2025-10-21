package services

import (
	"fmt"
	"net/http"

	"github.com/lapotkin/file-storage/internal/adapter/krakend"
	"github.com/lapotkin/file-storage/pkg/models"
)

type folderService struct {
	kc *krakend.Client
}

func NewFolderService(krakendClient *krakend.Client) FolderService {
	return &folderService{
		kc: krakendClient,
	}
}

func (s *folderService) CreateFolder(name, parentID string, token string) (*models.Folder, error) {
	if token == "" {
		return nil, ErrUnauthorized
	}

	reqBody := map[string]interface{}{
		"name":      name,
		"parent_id": parentID,
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodPost,
		Path:   "/api/v1/folders",
		Body:   reqBody,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	if err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to send request", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, NewAPIError(resp.StatusCode, "failed to create folder", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var folder models.Folder
	if err := resp.DecodeJSON(&folder); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return &folder, nil
}

func (s *folderService) GetFolder(folderID string, token string) (*models.Folder, error) {
	if folderID == "" {
		return nil, NewAPIError(http.StatusBadRequest, "folder ID is required", nil)
	}
	if token == "" {
		return nil, ErrUnauthorized
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/api/v1/folders/%s", folderID),
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
		return nil, NewAPIError(resp.StatusCode, "failed to get folder", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var folder models.Folder
	if err = resp.DecodeJSON(&folder); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return &folder, nil
}

func (s *folderService) ListFolders(parentID string, token string) (models.Folders, error) {
	if token == "" {
		return nil, ErrUnauthorized
	}

	path := "/api/v1/folders"
	if parentID != "" {
		path = fmt.Sprintf("%s?parent_id=%s", path, parentID)
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
		return nil, NewAPIError(resp.StatusCode, "failed to list folders", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var folders models.Folders
	if err = resp.DecodeJSON(&folders); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return folders, nil
}

func (s *folderService) UpdateFolder(folderID string, folder *models.Folder, token string) error {
	if folderID == "" {
		return NewAPIError(http.StatusBadRequest, "folder ID is required", nil)
	}
	if token == "" {
		return ErrUnauthorized
	}
	if folder == nil {
		return NewAPIError(http.StatusBadRequest, "folder data is required", nil)
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/api/v1/folders/%s", folderID),
		Body:   folder,
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
		return NewAPIError(resp.StatusCode, "failed to update folder", fmt.Errorf("response: %s", string(resp.Body)))
	}

	return nil
}

func (s *folderService) DeleteFolder(folderID string, token string) error {
	if folderID == "" {
		return NewAPIError(http.StatusBadRequest, "folder ID is required", nil)
	}
	if token == "" {
		return ErrUnauthorized
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/api/v1/folders/%s", folderID),
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
		return NewAPIError(resp.StatusCode, "failed to delete folder", fmt.Errorf("response: %s", string(resp.Body)))
	}

	return nil
}

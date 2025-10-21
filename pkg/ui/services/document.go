package services

import (
	"fmt"
	"net/http"

	"github.com/lapotkin/file-storage/internal/adapter/krakend"
	"github.com/lapotkin/file-storage/pkg/models"
)

type documentService struct {
	kc *krakend.Client
}

func NewDocumentService(kc *krakend.Client) DocumentService {
	return &documentService{
		kc: kc,
	}
}

func (s *documentService) CreateDocument(name, folderID string, token string) (*models.Document, error) {
	if token == "" {
		return nil, ErrUnauthorized
	}

	reqBody := map[string]interface{}{
		"name":      name,
		"folder_id": folderID,
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodPost,
		Path:   "/api/v1/documents",
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
		return nil, NewAPIError(resp.StatusCode, "failed to create document", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var document models.Document
	if err = resp.DecodeJSON(&document); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return &document, nil
}

func (s *documentService) GetDocument(documentID string, token string) (*models.Document, error) {
	if documentID == "" {
		return nil, NewAPIError(http.StatusBadRequest, "document ID is required", nil)
	}
	if token == "" {
		return nil, ErrUnauthorized
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/api/v1/documents/%s", documentID),
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
		return nil, NewAPIError(resp.StatusCode, "failed to get document", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var document models.Document
	if err = resp.DecodeJSON(&document); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return &document, nil
}

func (s *documentService) ListDocuments(folderID string, token string) (models.Documents, error) {
	if token == "" {
		return nil, ErrUnauthorized
	}

	path := "/api/v1/documents"
	if folderID != "" {
		path = fmt.Sprintf("%s?folder_id=%s", path, folderID)
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
		return nil, NewAPIError(resp.StatusCode, "failed to list documents", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var documents models.Documents
	if err = resp.DecodeJSON(&documents); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return documents, nil
}

func (s *documentService) UpdateDocument(documentID string, doc *models.Document, token string) error {
	if documentID == "" {
		return NewAPIError(http.StatusBadRequest, "document ID is required", nil)
	}
	if token == "" {
		return ErrUnauthorized
	}
	if doc == nil {
		return NewAPIError(http.StatusBadRequest, "document data is required", nil)
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/api/v1/documents/%s", documentID),
		Body:   doc,
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
		return NewAPIError(resp.StatusCode, "failed to update document", fmt.Errorf("response: %s", string(resp.Body)))
	}

	return nil
}

func (s *documentService) DeleteDocument(documentID string, token string) error {
	if documentID == "" {
		return NewAPIError(http.StatusBadRequest, "document ID is required", nil)
	}
	if token == "" {
		return ErrUnauthorized
	}

	resp, err := s.kc.Do(&krakend.Request{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/api/v1/documents/%s", documentID),
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
		return NewAPIError(resp.StatusCode, "failed to delete document", fmt.Errorf("response: %s", string(resp.Body)))
	}

	return nil
}

func (s *documentService) SearchDocuments(query string, token string) (models.Documents, error) {
	if token == "" {
		return nil, ErrUnauthorized
	}

	path := "/api/v1/documents/search"
	if query != "" {
		path = fmt.Sprintf("%s?q=%s", path, query)
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
		return nil, NewAPIError(resp.StatusCode, "failed to search documents", fmt.Errorf("response: %s", string(resp.Body)))
	}

	var documents models.Documents
	if err = resp.DecodeJSON(&documents); err != nil {
		return nil, NewAPIError(http.StatusInternalServerError, "failed to decode response", err)
	}

	return documents, nil
}

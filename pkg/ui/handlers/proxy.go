package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lapotkin/file-storage/internal/adapter/krakend"
)

func NewProxyHandler(krakendClient *krakend.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		if query != "" {
			path += "?" + query
		}

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		token, _ := c.Cookie("token")

		headers := make(map[string]string)

		if contentType := c.GetHeader("Content-Type"); contentType != "" {
			headers["Content-Type"] = contentType
		}
		if accept := c.GetHeader("Accept"); accept != "" {
			headers["Accept"] = accept
		}

		if token != "" {
			headers["Authorization"] = "Bearer " + token
		}

		resp, err := krakendClient.Do(&krakend.Request{
			Method:  c.Request.Method,
			Path:    path,
			Body:    bodyBytes,
			Headers: headers,
		})

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "failed to forward request to API gateway",
			})
			return
		}

		for key, values := range resp.Headers {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		c.Data(resp.StatusCode, resp.ContentType(), resp.Body)
	}
}

func NewAPIProxyHandler(krakendClient *krakend.Client, _ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		if query != "" {
			path += "?" + query
		}

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		token, _ := c.Cookie("token")

		headers := make(map[string]string)

		for key, values := range c.Request.Header {
			if len(values) > 0 && !isHopByHopHeader(key) {
				headers[key] = values[0]
			}
		}

		if token != "" {
			headers["Authorization"] = "Bearer " + token
		}

		resp, err := krakendClient.Do(&krakend.Request{
			Method:  c.Request.Method,
			Path:    path,
			Body:    bodyBytes,
			Headers: headers,
		})

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "failed to forward request to API gateway",
			})
			return
		}

		for key, values := range resp.Headers {
			if !isHopByHopHeader(key) {
				for _, value := range values {
					c.Header(key, value)
				}
			}
		}

		c.Data(resp.StatusCode, resp.ContentType(), resp.Body)
	}
}

func isHopByHopHeader(header string) bool {
	hopByHopHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	header = strings.ToLower(header)
	for _, h := range hopByHopHeaders {
		if strings.ToLower(h) == header {
			return true
		}
	}
	return false
}

func NewDocumentUploadProxyHandler(krakendClient *krakend.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}
		defer c.Request.Body.Close()

		contentType := c.GetHeader("Content-Type")

		req, err := http.NewRequest(http.MethodPost, krakendClient.BaseURL()+"/api/v1/documents", bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to upload document: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
			return
		}

		c.Data(resp.StatusCode, "application/json", respBody)
	}
}

func NewDocumentDownloadProxyHandler(krakendClient *krakend.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		documentID := c.Param("id")
		if documentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "document ID is required"})
			return
		}

		resp, err := krakendClient.Do(&krakend.Request{
			Method: http.MethodGet,
			Path:   "/api/v1/documents/" + documentID + "/download",
			Headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
		})

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to download document"})
			return
		}

		if filename := resp.Headers.Get("Content-Disposition"); filename != "" {
			c.Header("Content-Disposition", filename)
		}
		if contentType := resp.ContentType(); contentType != "" {
			c.Header("Content-Type", contentType)
		}

		c.Data(resp.StatusCode, resp.ContentType(), resp.Body)
	}
}

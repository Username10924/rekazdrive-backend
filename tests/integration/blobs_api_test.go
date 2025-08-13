package integration

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlobsAPI_LocalFS(t *testing.T) {
    baseURL := "http://localhost:8080/v1"
    
    // Login to get token
    token := login(t, baseURL)
    
    // Create blob
    blobID := "test-id"
    data := []byte("hello world :)")
    createBlob(t, baseURL, token, blobID, data)
    
    // Get blob and verify
    blob := getBlob(t, baseURL, token, blobID)
    require.Equal(t, blobID, blob["id"])
}

func login(t *testing.T, baseURL string) string {
    body := map[string]string{
        "username": "admin",
        "password": "admin",
    }
    
    jsonBody, _ := json.Marshal(body) // map to JSON
    req, _ := http.NewRequest("POST", baseURL+"/auth/login", bytes.NewReader(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)
    
    var result map[string]any
    json.NewDecoder(resp.Body).Decode(&result) // JSON response to map
    
    return result["token"].(string)
}

func createBlob(t *testing.T, baseURL, token, blobID string, data []byte) {
    body := map[string]string{
        "id":   blobID,
        "data": base64.StdEncoding.EncodeToString(data),
    }
    
    jsonBody, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", baseURL+"/blobs", bytes.NewReader(jsonBody))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    require.Equal(t, http.StatusCreated, resp.StatusCode)
}

func getBlob(t *testing.T, baseURL, token, blobID string) map[string]any {
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/blobs/%s", baseURL, blobID), nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)
    
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    return result
}
package unit

import (
	"os"
	"rekazdrive/internal/storage"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestS3Backend_SaveLoad(t *testing.T) {
    // load .env file for testing (same as main.go does)
    _ = godotenv.Load("../../.env")
    
    endpoint := os.Getenv("S3_ENDPOINT")
    bucket := os.Getenv("S3_BUCKET")
    accessKey := os.Getenv("S3_ACCESS_KEY")
    secret := os.Getenv("S3_SECRET")
    region := os.Getenv("S3_REGION")

	backend := storage.NewS3Backend(endpoint, bucket, accessKey, secret, region)

	data := []byte("hello world :)")
	id := "test-id"

	err := backend.Save(id, data)
	require.NoError(t, err)
	
	loadedData, err := backend.Load(id)
	require.NoError(t, err)
	require.Equal(t, data, loadedData)

	// Clean up
	err = backend.Delete(id)
	require.NoError(t, err)
}

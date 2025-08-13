package unit

import (
	"rekazdrive/internal/storage"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalBackend_SaveLoad(t *testing.T) {
	tmpDir := t.TempDir()
	backend := storage.NewLocalBackend(tmpDir)

	data := []byte("hello world :)")
	id := "test-id"

	err := backend.Save(id, data)
	require.NoError(t, err)
	
	loadedData, err := backend.Load(id)
	require.NoError(t, err)
	require.Equal(t, data, loadedData)
}
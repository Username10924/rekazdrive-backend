package unit

import (
	"database/sql"
	"rekazdrive/internal/db"
	"rekazdrive/internal/storage"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/require"
)

func TestDBBlobBackend_SaveLoad(t *testing.T) {
    dsn := "postgres://postgres:MxWU9s3aO25V@localhost:5432/blob_data?sslmode=disable"

    database, err := sql.Open("postgres", dsn)
    require.NoError(t, err)
    defer database.Close()

    // Initialize blob table
    err = db.InitBlobTable(database)
    require.NoError(t, err)

    backend := storage.NewDBBlobBackend(database)

    data := []byte("hello world :)")
    id := "test-id"

    err = backend.Save(id, data)
    require.NoError(t, err)
    
    loadedData, err := backend.Load(id)
    require.NoError(t, err)
    require.Equal(t, data, loadedData)
}
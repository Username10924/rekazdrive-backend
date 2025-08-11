package db

import (
	"database/sql"
	"time"
)

type MetadataDB struct {
	DB *sql.DB
}

type BlobMeta struct {
	ID string
	Size int
	CreatedAt time.Time
}

// initialize new MetadataDB using Postgres
func NewPostgres(dsn string) (*MetadataDB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &MetadataDB{DB: db}, nil
}

func (m *MetadataDB) InitSchema() error {
	query := `CREATE TABLE IF NOT EXISTS blobs_metadata (
		id TEXT PRIMARY KEY,
		size INTEGER NOT NULL,
		created_at TIMESTAMPTZ NOT NULL
		);`
	
	_, err := m.DB.Exec(query)

	return err
}

func (m *MetadataDB) SaveMetadata(id string, size int, createdAt time.Time) error {
	query := `INSERT INTO blobs_metadata(id, size, created_at)
		      VALUES($1, $2, $3)
			  ON CONFLICT (id) DO UPDATE
			  SET size = EXCLUDED.size, created_at = EXCLUDED.created_at;`
	
	_, err := m.DB.Exec(query, id, size, createdAt.UTC())

	return err
}

func (m *MetadataDB) GetMetadata(id string) (*BlobMeta, error) {
	query := `SELECT id, size, created_at FROM blobs_metadata WHERE id = $1;`
	row := m.DB.QueryRow(query, id)

	var meta BlobMeta
	err := row.Scan(&meta.ID, &meta.Size, &meta.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}
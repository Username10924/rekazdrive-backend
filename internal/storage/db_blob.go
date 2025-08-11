package storage

import (
	"database/sql"
	"time"
)

type DBBlobBackend struct {
	db *sql.DB
}

func NewDBBlobBackend(db *sql.DB) *DBBlobBackend {
	return &DBBlobBackend{db: db}
}

func (d *DBBlobBackend) Save(id string, data []byte) error {
	query := `INSERT OR REPLACE INTO blobs_data(id, data, created_at) VALUES ($1, $2, $3)`
	_, err := d.db.Exec(query, id, data, time.Now().UTC().Format(time.RFC3339))
	return err
}

func (d *DBBlobBackend) Load(id string) ([]byte, error) {
	query := `SELECT data FROM blobs_data WHERE id = $1`
	row := d.db.QueryRow(query, id)
	var data []byte
	err := row.Scan(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *DBBlobBackend) Delete(id string) error {
	query := `DELETE FROM blobs_data WHERE id = $1`
	_, err := d.db.Exec(query, id)
	return err
}
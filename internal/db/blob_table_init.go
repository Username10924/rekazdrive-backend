package db

import "database/sql"

func InitBlobTable(sqlDB *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS blobs_data (
		id TEXT PRIMARY KEY,
		data BYTEA NOT NULL,
		created_at TIMESTAMPTZ NOT NULL
		);`
	_, err := sqlDB.Exec(query)

	return err
}
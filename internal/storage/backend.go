package storage

// abstract interface to be implemented by different storage methods
type StorageBackend interface {
	Save(id string, data []byte) error
	Load(id string) ([]byte, error)
	Delete (id string) error
}
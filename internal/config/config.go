package config

import "os"

type Config struct {
	AuthToken string // for API access
	StorageBackend string

	// Local
	LocalPath string

	// S3
	S3Endpoint string
	S3Bucket string
	S3AccessKey string
	S3SecretKey string
	S3Region string

	// FTP
	FTPHost string
	FTPUser string
	FTPPass string

	// DB DSN
	MetadataDSN string
	BlobDBDSN string
}

func LoadFromEnv() Config {
	return Config{
		AuthToken : os.Getenv("AUTH_TOKEN"),
		StorageBackend: os.Getenv("STORAGE_BACKEND"),
		LocalPath: os.Getenv("LOCAL_PATH"),
		S3Endpoint: os.Getenv("S3_ENDPOINT"),
		S3Bucket: os.Getenv("S3_BUCKET"),
		S3AccessKey: os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey: os.Getenv("S3_SECRET_KEY"),
		S3Region: os.Getenv("S3_REGION"),
		FTPHost: os.Getenv("FTP_HOST"),
		FTPUser: os.Getenv("FTP_USER"),
		FTPPass: os.Getenv("FTP_PASS"),
		MetadataDSN: os.Getenv("METADATA_DSN"),
		BlobDBDSN: os.Getenv("BLOB_DB_DSN"),
	}
}
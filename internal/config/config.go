package config

import "os"

type Config struct {
	JWTSecret string
	JWTExpiration string

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
	FTPBasePath string

	// DB DSN
	MetadataDSN string
	BlobDBDSN string

	// Auth credentials
	AdminUser string
	AdminPass string
}

func LoadFromEnv() Config {
	return Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTExpiration: os.Getenv("JWT_EXPIRATION"),
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
		FTPBasePath: os.Getenv("FTP_BASE_PATH"),
		MetadataDSN: os.Getenv("METADATA_DSN"),
		BlobDBDSN: os.Getenv("BLOB_DB_DSN"),
		AdminUser: os.Getenv("ADMIN_USER"),
		AdminPass: os.Getenv("ADMIN_PASS"),
	}
}
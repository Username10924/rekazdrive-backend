package server

import (
	"database/sql"
	"log"
	"rekazdrive/internal/config"
	"rekazdrive/internal/db"
	"rekazdrive/internal/handlers"
	"rekazdrive/internal/middleware"
	"rekazdrive/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadFromEnv()

	metaDB, err := db.NewPostgres(cfg.MetadataDSN)
	if err != nil {
		log.Fatalf("Failed to connect to metadata database: %v", err)
	}
	if err := metaDB.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize metadata schema: %v", err)
	}

	// storage backend selection from env
	var backend storage.StorageBackend
	switch cfg.StorageBackend {
	case "local":
		backend = storage.NewLocalBackend(cfg.LocalPath)
	case "db":
		blobDB, err := sql.Open("postgres", cfg.BlobDBDSN)
		if err != nil {
			log.Fatalf("Failed to connect to blob database: %v", err)
		}
		if err := db.InitBlobTable(blobDB); err != nil {
			log.Fatalf("Failed to initialize blob table: %v", err)
		}
		backend = storage.NewDBBlobBackend(blobDB)
	case "ftp":
		backend = storage.NewFTPBackend(cfg.FTPHost, cfg.FTPUser, cfg.FTPPass, cfg.FTPBasePath)
	case "s3":
		backend = storage.NewS3Backend(cfg.S3Endpoint, cfg.S3Bucket, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Region)
	default:
		log.Fatalf("Unsupported storage backend: %s", cfg.StorageBackend)
	}

	router := gin.Default()

	// public group
	v1 := router.Group("/v1")
	v1.POST("/auth/login", handlers.LoginHandler(cfg))

	// protected group
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	blobs := protected.Group("/blobs")
	{
		blobs.POST("", handlers.NewBlobHandler(backend, metaDB).PostBlob)
		blobs.GET("/:id", handlers.NewBlobHandler(backend, metaDB).GetBlob)
	}

	port := "8080"
	log.Printf("Starting server on port %s", port)
	router.Run(":" + port)
}
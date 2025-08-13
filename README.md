# RekAZ Drive Backend

Blob storage API with JWT auth and multiple storage backends (local, S3, FTP, PostgreSQL).

## Prerequisites

- Go 1.22+
- Docker & Docker Compose
- PostgreSQL

## Setup

```bash
# Clone and start
git clone https://github.com/Username10924/rekazdrive-backend.git
cd rekazdrive-backend

# Run with Docker
docker-compose up -d

# Or run locally
go run main.go
```

## Usage

```bash
# Login
curl -X POST localhost:8080/v1/auth/login \
  -d '{"username":"admin","password":"admin"}'

# Store blob
curl -X POST localhost:8080/v1/blobs \
  -H "Authorization: Bearer TOKEN" \
  -d '{"id":"test","data":"base64data"}'

# Get blob
curl localhost:8080/v1/blobs/test \
  -H "Authorization: Bearer TOKEN"
```

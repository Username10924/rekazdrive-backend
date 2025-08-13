package handlers

import (
	"encoding/base64"
	"net/http"
	"rekazdrive/internal/db"
	"rekazdrive/internal/storage"
	"time"

	"github.com/gin-gonic/gin"
)

type BlobHandler struct {
	Store storage.StorageBackend
	Meta *db.MetadataDB
}

func NewBlobHandler(store storage.StorageBackend, meta *db.MetadataDB) *BlobHandler {
	return &BlobHandler{
		Store: store,
		Meta:  meta,
	}
}

type postReq struct {
	ID string `json:"id"`
	Data string `json:"data"`
}

type getResp struct {
	ID string `json:"id"`
	Data string `json:"data"` // Base64 encoded data
	Size int `json:"size"`
	CreatedAt string `json:"created_at"`
}

func (h *BlobHandler) PostBlob (c *gin.Context) {
	var r postReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if r.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	data, err := base64.StdEncoding.DecodeString(r.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid base64"})
		return
	}

	if err := h.Store.Save(r.ID, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed", "detail": err.Error()})
		return
	}

	if err := h.Meta.SaveMetadata(r.ID, len(data), time.Now().UTC()); err != nil {
		_ = h.Store.Delete(r.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "meta save failed", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (h *BlobHandler) GetBlob(c *gin.Context) {
	id := c.Param("id")
	meta, err := h.Meta.GetMetadata(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	data, err := h.Store.Load(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	resp := getResp{
		ID:        id,
		Data:      base64.StdEncoding.EncodeToString(data),
		Size:      meta.Size,
		CreatedAt: meta.CreatedAt.UTC().Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, resp)
}
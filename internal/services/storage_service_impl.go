package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/ports/services"
	"github.com/gildo-cordeiro/mapleplan-api/pkg/utils"
	"github.com/google/uuid"
)

type StorageServiceImpl struct {
	storage services.StorageService
}

func NewStorageService(storage services.StorageService) *StorageServiceImpl {
	return &StorageServiceImpl{
		storage: storage,
	}
}

// UploadFile uploads a file to storage with a unique key
func (s *StorageServiceImpl) UploadFile(ctx context.Context, bucket string, file io.Reader, originalFilename string, fileSize int64, contentType string) (string, error) {
	if bucket == "" {
		return "", fmt.Errorf("UploadFile called with empty bucket name")
	}

	if fileSize == 0 || fileSize > 100*1024*1024 { // 100MB limit
		return "", fmt.Errorf("invalid file size: %d bytes", fileSize)
	}

	// Generate unique file key with timestamp
	fileExt := filepath.Ext(originalFilename)
	uniqueKey := fmt.Sprintf("uploads/%s%s", uuid.New().String(), fileExt)

	utils.Log.Logger.Infof("Starting upload for file: %s (size: %d bytes)", originalFilename, fileSize)

	url, err := s.storage.UploadFile(ctx, bucket, uniqueKey, file, fileSize, contentType)
	if err != nil {
		utils.Log.Logger.Errorf("Failed to upload file: %v", err)
		return "", err
	}

	utils.Log.Logger.Infof("File uploaded successfully: %s -> %s", originalFilename, uniqueKey)
	return url, nil
}

// DeleteFile removes a file from storage
func (s *StorageServiceImpl) DeleteFile(ctx context.Context, bucket, fileKey string) error {
	if bucket == "" || fileKey == "" {
		return fmt.Errorf("bucket and file key are required")
	}

	utils.Log.Logger.Infof("Deleting file: %s/%s", bucket, fileKey)
	return s.storage.DeleteFile(ctx, bucket, fileKey)
}

// GetFileURL returns a presigned URL for accessing the file
func (s *StorageServiceImpl) GetFileURL(ctx context.Context, bucket, fileKey string) (string, error) {
	if bucket == "" || fileKey == "" {
		return "", fmt.Errorf("bucket and file key are required")
	}

	return s.storage.GetFileURL(ctx, bucket, fileKey)
}

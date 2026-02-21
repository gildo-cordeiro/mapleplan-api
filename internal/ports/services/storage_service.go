package services

import (
	"context"
	"io"
)

// StorageService defines the interface for file storage operations
type StorageService interface {
	// UploadFile uploads a file to storage and returns the file URL
	UploadFile(ctx context.Context, bucket, key string, reader io.Reader, fileSize int64, contentType string) (string, error)

	// DeleteFile removes a file from storage
	DeleteFile(ctx context.Context, bucket, key string) error

	// GetFileURL returns the URL to access a file
	GetFileURL(ctx context.Context, bucket, key string) (string, error)

	// DownloadFile downloads a file from storage
	DownloadFile(ctx context.Context, bucket, key string) (io.Reader, error)

	// ListFiles lists all files in a bucket
	ListFiles(ctx context.Context, bucket, prefix string) ([]string, error)
}

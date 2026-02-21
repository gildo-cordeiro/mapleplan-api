package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sirupsen/logrus"
)

type S3Storage struct {
	client   *s3.Client
	uploader *manager.Uploader
	region   string
	endpoint string
	logger   *logrus.Logger
}

func NewS3Storage(client *s3.Client, region, endpoint string, logger *logrus.Logger) *S3Storage {
	return &S3Storage{
		client:   client,
		uploader: manager.NewUploader(client),
		region:   region,
		endpoint: endpoint,
		logger:   logger,
	}
}

func (s *S3Storage) UploadFile(ctx context.Context, bucket, key string, reader io.Reader, fileSize int64, contentType string) (string, error) {
	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPrivate,
	})
	if err != nil {
		s.logger.Errorf("Failed to upload file to S3: %v", err)
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	s.logger.Infof("File uploaded successfully to S3: %s/%s", bucket, key)
	return fmt.Sprintf("s3://%s/%s", bucket, key), nil
}

func (s *S3Storage) DeleteFile(ctx context.Context, bucket, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.Errorf("Failed to delete file from S3: %v", err)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	s.logger.Infof("File deleted successfully from S3: %s/%s", bucket, key)
	return nil
}

func (s *S3Storage) GetFileURL(ctx context.Context, bucket, key string) (string, error) {
	if s.endpoint != "" {
		base := strings.TrimRight(s.endpoint, "/")
		url := fmt.Sprintf("%s/%s/%s", base, bucket, key)
		s.logger.Infof("Generated S3 URL: %s", url)
		return url, nil
	}

	// For production, consider using CloudFront or Signed URLs through AWS SDK.
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, s.region, key)
	s.logger.Infof("Generated S3 URL: %s", url)
	return url, nil
}

func (s *S3Storage) DownloadFile(ctx context.Context, bucket, key string) (io.Reader, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.Errorf("Failed to download file from S3: %v", err)
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return result.Body, nil
}

func (s *S3Storage) ListFiles(ctx context.Context, bucket, prefix string) ([]string, error) {
	var files []string

	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			s.logger.Errorf("Error listing objects: %v", err)
			return nil, fmt.Errorf("error listing objects: %w", err)
		}

		for _, obj := range page.Contents {
			files = append(files, *obj.Key)
		}
	}

	return files, nil
}

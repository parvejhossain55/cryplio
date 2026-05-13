package storage

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"strings"

	"cryplio/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// UploadInput describes a file upload request.
type UploadInput struct {
	Key         string
	ContentType string
	Body        []byte
}

// UploadResult describes the stored file location.
type UploadResult struct {
	Key string
	URL string
}

// ObjectStorage handles evidence and document uploads.
type ObjectStorage interface {
	Upload(ctx context.Context, input UploadInput) (*UploadResult, error)
	Delete(ctx context.Context, key string) error
}

// S3Storage implements ObjectStorage using MinIO/S3-compatible storage.
type S3Storage struct {
	client        *minio.Client
	endpoint      string
	accessKey     string
	secretKey     string
	useSSL        bool
	bucket        string
	publicBaseURL string // Optional public URL for accessing objects
}

// NewS3Storage creates a new S3Storage instance using configuration.
func NewS3Storage(cfg *config.Config) (*S3Storage, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKeyID, cfg.S3SecretAccessKey, ""),
		Secure: cfg.S3UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Create bucket if it doesn't exist
	exists, err := minioClient.BucketExists(context.Background(), cfg.S3BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), cfg.S3BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Set bucket policy to allow public reads for objects
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, cfg.S3BucketName)
	if err := minioClient.SetBucketPolicy(context.Background(), cfg.S3BucketName, policy); err != nil {
		// Log but don't fail startup if policy cannot be set
		fmt.Printf("Warning: failed to set bucket policy: %v\n", err)
	}

	return &S3Storage{
		client:        minioClient,
		endpoint:      cfg.S3Endpoint,
		accessKey:     cfg.S3AccessKeyID,
		secretKey:     cfg.S3SecretAccessKey,
		useSSL:        cfg.S3UseSSL,
		bucket:        cfg.S3BucketName,
		publicBaseURL: cfg.S3PublicBaseURL,
	}, nil
}

// Upload stores a file in S3-compatible storage.
func (s *S3Storage) Upload(ctx context.Context, input UploadInput) (*UploadResult, error) {
	if input.ContentType == "" {
		input.ContentType = mime.TypeByExtension(strings.ToLower(input.Key))
	}
	if input.ContentType == "" {
		input.ContentType = "application/octet-stream"
	}

	reader := bytes.NewReader(input.Body)
	_, err := s.client.PutObject(ctx, s.bucket, input.Key, reader, int64(len(input.Body)), minio.PutObjectOptions{
		ContentType: input.ContentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload object: %w", err)
	}

	// Construct public URL
	var url string
	if s.publicBaseURL != "" {
		url = fmt.Sprintf("%s/%s", strings.TrimRight(s.publicBaseURL, "/"), input.Key)
	} else {
		schema := "http"
		if s.useSSL {
			schema = "https"
		}
		url = fmt.Sprintf("%s://%s/%s/%s", schema, s.endpoint, s.bucket, input.Key)
	}

	return &UploadResult{
		Key: input.Key,
		URL: url,
	}, nil
}

// Delete removes a file from S3-compatible storage.
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

package storage

import "context"

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

// NoopS3Storage is a placeholder until S3-compatible object storage is integrated.
type NoopS3Storage struct{}

func NewS3Storage() *NoopS3Storage {
	return &NoopS3Storage{}
}

func (s *NoopS3Storage) Upload(context.Context, UploadInput) (*UploadResult, error) {
	return nil, nil
}

func (s *NoopS3Storage) Delete(context.Context, string) error {
	return nil
}

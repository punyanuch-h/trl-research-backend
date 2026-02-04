package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

type GCSClient struct {
	BucketName string
}

func NewGCSClient(bucketName string) (*GCSClient, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name is required")
	}
	return &GCSClient{
		BucketName: bucketName,
	}, nil
}

// Direct Upload (BACKEND → GCS)
func (c *GCSClient) UploadFile(
	objectPath string,
	contentType string,
	file io.Reader,
) error {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	wc := client.Bucket(c.BucketName).Object(objectPath).NewWriter(ctx)
	wc.ContentType = contentType

	if _, err = io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	return nil
}

// Generate Upload Signed URL
// ✅ IAM / ADC compatible
func (c *GCSClient) GenerateUploadSignedURL(
	objectPath string,
	expireMinutes int,
) (string, error) {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Method:  "PUT",
		Expires: time.Now().Add(time.Duration(expireMinutes) * time.Minute),
		Scheme:  storage.SigningSchemeV4,
	}

	url, err := client.
		Bucket(c.BucketName).
		SignedURL(objectPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload signed url: %w", err)
	}

	return url, nil
}

// Generate Download Signed URL
// ✅ IAM / ADC compatible
func (c *GCSClient) GenerateDownloadSignedURL(
	objectPath string,
	expireMinutes int,
) (string, error) {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(expireMinutes) * time.Minute),
		Scheme:  storage.SigningSchemeV4,
	}

	url, err := client.
		Bucket(c.BucketName).
		SignedURL(objectPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate download signed url: %w", err)
	}

	return url, nil
}

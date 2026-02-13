package storage

import (
	"context"
	"fmt"
	"io"
	"log"
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
	contentType string,
	expireMinutes int,
) (string, error) {
	log.Printf("[GCSClient] Generating Upload Signed URL for bucket: %s, object: %s, contentType: %s\n", c.BucketName, objectPath, contentType)

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("[GCSClient] Error creating storage client: %v\n", err)
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Method:      "PUT",
		ContentType: contentType,
		Expires:     time.Now().Add(time.Duration(expireMinutes) * time.Minute),
		Scheme:      storage.SigningSchemeV4,
	}

	log.Printf("[GCSClient] SignedURL Options: Method=%s, ContentType=%s, Expires=%v, Scheme=V4\n", opts.Method, opts.ContentType, opts.Expires)

	url, err := client.
		Bucket(c.BucketName).
		SignedURL(objectPath, opts)
	if err != nil {
		log.Printf("[GCSClient] Error from SignedURL: %v\n", err)
		return "", fmt.Errorf("failed to generate upload signed url: %w", err)
	}

	log.Printf("[GCSClient] Successfully generated Signed URL length: %d\n", len(url))
	return url, nil
}

// Generate Download Signed URL
// ✅ IAM / ADC compatible
func (c *GCSClient) GenerateDownloadSignedURL(
	objectPath string,
	expireMinutes int,
) (string, error) {
	log.Printf("[GCSClient] Generating Download Signed URL for bucket: %s, object: %s\n", c.BucketName, objectPath)

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("[GCSClient] Error creating storage client: %v\n", err)
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(expireMinutes) * time.Minute),
		Scheme:  storage.SigningSchemeV4,
	}

	log.Printf("[GCSClient] SignedURL Options: Method=%s, Expires=%v, Scheme=V4\n", opts.Method, opts.Expires)

	url, err := client.
		Bucket(c.BucketName).
		SignedURL(objectPath, opts)
	if err != nil {
		log.Printf("[GCSClient] Error from SignedURL: %v\n", err)
		return "", fmt.Errorf("failed to generate download signed url: %w", err)
	}

	log.Printf("[GCSClient] Successfully generated Signed URL length: %d\n", len(url))
	return url, nil
}

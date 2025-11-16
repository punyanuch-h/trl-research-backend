package storage

import (
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type GCSClient struct {
	BucketName string
	SAEmail    string
}

// NewGCSClient returns a new client struct (no actual network call)
func NewGCSClient(bucketName, serviceAccountEmail string) *GCSClient {
	return &GCSClient{
		BucketName: bucketName,
		SAEmail:    serviceAccountEmail,
	}
}

// GenerateUploadSignedURL creates a signed URL for uploading a file (PUT)
func (c *GCSClient) GenerateUploadSignedURL(objectPath, contentType string, expireMinutes int) (string, error) {

	// path to JSON file
	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if creds == "" {
		return "", fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	// load private key from JSON
	privateKey, err := loadPrivateKeyFromCredentials(creds)
	if err != nil {
		return "", fmt.Errorf("loading private key: %w", err)
	}

	// set signed URL options
	opts := &storage.SignedURLOptions{
		Method:         "PUT",
		Expires:        time.Now().Add(time.Duration(expireMinutes) * time.Minute),
		Scheme:         storage.SigningSchemeV4,
		GoogleAccessID: c.SAEmail,
		PrivateKey:     privateKey,
		ContentType:    contentType,
	}

	// create signed url
	url, err := storage.SignedURL(c.BucketName, objectPath, opts)
	if err != nil {
		return "", fmt.Errorf("cannot generate signed url: %w", err)
	}

	return url, nil
}

func (c *GCSClient) GenerateDownloadSignedURL(objectPath string, expireMinutes int) (string, error) {
    
    creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
    if creds == "" {
        return "", fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS not set")
    }

    privateKey, err := loadPrivateKeyFromCredentials(creds)
    if err != nil {
        return "", fmt.Errorf("failed to load private key: %w", err)
    }

    opts := &storage.SignedURLOptions{
        Method:         "GET",
        Expires:        time.Now().Add(time.Duration(expireMinutes) * time.Minute),
        Scheme:         storage.SigningSchemeV4,
        GoogleAccessID: c.SAEmail,
        PrivateKey:     privateKey,
    }

    url, err := storage.SignedURL(c.BucketName, objectPath, opts)
    if err != nil {
        return "", fmt.Errorf("failed to generate signed download URL: %w", err)
    }

    return url, nil
}

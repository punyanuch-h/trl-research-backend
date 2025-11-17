package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type serviceAccountKey struct {
	PrivateKey string `json:"private_key"`
}

// loadPrivateKeyFromCredentials loads a private key from a service account JSON key file
func loadPrivateKeyFromCredentials(filePath string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read service account file: %w", err)
	}

	var key serviceAccountKey
	if err := json.Unmarshal(fileBytes, &key); err != nil {
		return nil, fmt.Errorf("cannot unmarshal json: %w", err)
	}

	if key.PrivateKey == "" {
		return nil, fmt.Errorf("private_key field is empty in JSON file")
	}

	return []byte(key.PrivateKey), nil
}

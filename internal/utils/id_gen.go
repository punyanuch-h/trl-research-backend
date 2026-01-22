package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenerateID generates a collision-free ID with a given prefix.
// The ID format is <prefix>-<uuid_segment>
func GenerateID(prefix string) string {
	newUUID := uuid.New().String()
	// Using the full UUID to ensure uniqueness as requested
	return fmt.Sprintf("%s-%s", prefix, newUUID)
}

// GenerateShortID generates a shorter collision-free ID with a given prefix.
// It uses the first segment of the UUID.
func GenerateShortID(prefix string) string {
	newUUID := uuid.New().String()
	segment := strings.Split(newUUID, "-")[0]
	return fmt.Sprintf("%s-%s", prefix, segment)
}

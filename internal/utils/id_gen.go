package utils

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// GenerateID generates an ID with the format PREFIX-00001
// It finds the highest numeric ID with the given prefix and increments it.
func GenerateID(db *gorm.DB, tableName string, prefix string) (string, error) {
	var ids []string
	// Fetch IDs that start with the prefix
	err := db.Table(tableName).Where("id LIKE ?", prefix+"-%").Select("id").Find(&ids).Error
	if err != nil {
		return "", err
	}

	maxNum := 0
	for _, id := range ids {
		// Example: CS-00001
		parts := strings.Split(id, "-")
		if len(parts) != 2 {
			continue
		}

		num, err := strconv.Atoi(parts[1])
		if err == nil {
			if num > maxNum {
				maxNum = num
			}
		}
	}

	return fmt.Sprintf("%s-%05d", prefix, maxNum+1), nil
}

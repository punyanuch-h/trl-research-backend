package utils

import (
	"strings"
)

// AttachmentKeys defines the allowed semantic keys and their corresponding entity names
var AttachmentKeys = map[string]string{
	"cases_attachments":       "cases",
	"assessments_attachments": "assessments",
	"ips_attachments":         "ips",
}

// ExtractAttachments finds all keys ending with "_attachments" in the input map,
// validates them against the allow-list, and returns a map of entity name to paths.
func ExtractAttachments(body map[string]interface{}) map[string][]string {
	result := make(map[string][]string)

	for key, value := range body {
		if !strings.HasSuffix(key, "_attachments") {
			continue
		}

		entity, allowed := AttachmentKeys[key]
		if !allowed {
			continue
		}

		// Validate value is []interface{} (which it will be if it's a JSON array) or []string
		var paths []string
		switch v := value.(type) {
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok && s != "" {
					paths = append(paths, s)
				}
			}
		case []string:
			for _, s := range v {
				if s != "" {
					paths = append(paths, s)
				}
			}
		case string:
			// Handle single string if frontend sends it that way, but requirement says "array of strings"
			if v != "" {
				paths = append(paths, v)
			}
		}

		if len(paths) > 0 {
			result[entity] = paths
		}
	}

	return result
}

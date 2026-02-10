package utils

import (
	"reflect"
	"testing"
)

func TestExtractAttachments(t *testing.T) {
	tests := []struct {
		name     string
		body     map[string]interface{}
		expected map[string][]string
	}{
		{
			name: "Valid attachments",
			body: map[string]interface{}{
				"cases_attachments": []interface{}{"file1.pdf", "file2.jpg"},
				"ips_attachments":   []string{"file3.png"},
				"other_field":       "value",
			},
			expected: map[string][]string{
				"cases": {"file1.pdf", "file2.jpg"},
				"ips":   {"file3.png"},
			},
		},
		{
			name: "Invalid keys or types",
			body: map[string]interface{}{
				"unknown_attachments": []interface{}{"file1.pdf"},
				"cases_attachments":   123,
			},
			expected: map[string][]string{},
		},
		{
			name: "Empty values",
			body: map[string]interface{}{
				"cases_attachments": []interface{}{"", ""},
			},
			expected: map[string][]string{},
		},
		{
			name: "Single string value",
			body: map[string]interface{}{
				"cases_attachments": "file1.pdf",
			},
			expected: map[string][]string{
				"cases": {"file1.pdf"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractAttachments(tt.body)
			if !reflect.DeepEqual(got, tt.expected) {
				if len(got) == 0 && len(tt.expected) == 0 {
					return
				}
				t.Errorf("ExtractAttachments() = %v, want %v", got, tt.expected)
			}
		})
	}
}

//go:build xattr

package tool

import (
	"os"
	"testing"

	"github.com/pkg/xattr"
)

func TestStripCredentials(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "url with user and pass",
			input:    "https://user:pass@example.com/path",
			expected: "https://example.com/path",
		},
		{
			name:     "url with user only",
			input:    "ftp://user@ftp.example.com",
			expected: "ftp://ftp.example.com",
		},
		{
			name:     "url with no credentials",
			input:    "http://example.com",
			expected: "http://example.com",
		},
		{
			name:    "invalid url",
			input:   "://",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := stripCredentials(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("stripCredentials() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && result != tc.expected {
				t.Errorf("stripCredentials() = %q; want %q", result, tc.expected)
			}
		})
	}
}

func TestWriteXattr(t *testing.T) {
	// Create a temporary file to set attributes on.
	tempFile, err := os.CreateTemp(t.TempDir(), "xattr-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	filePath := tempFile.Name()
	tempFile.Close()

	// Sample data to write as xattrs.
	data := map[string]interface{}{
		"content_type":  "application/json",
		"referer":       "http://previous.example.com",
		"url_effective": "https://user:pass@current.example.com",
	}

	// Call the function to write the attributes.
	if err := WriteXattr(filePath, data); err != nil {
		// xattrs are not supported on all filesystems (e.g., tmpfs on some CIs).
		// If the error indicates this, we skip the test.
		if e, ok := err.(*xattr.Error); ok && e.Err.Error() == "operation not supported" {
			t.Skipf("Skipping xattr test: filesystem does not support extended attributes (%v)", err)
		}
		t.Fatalf("WriteXattr() failed: %v", err)
	}

	// Verify the attributes were written correctly.
	testCases := []struct {
		attrName string
		expected string
	}{
		{"user.creator", "curl-translation-go"},
		{"user.mime_type", "application/json"},
		{"user.xdg.referrer.url", "http://previous.example.com"},
		{"user.xdg.origin.url", "https://current.example.com"}, // Credentials should be stripped
	}

	for _, tc := range testCases {
		t.Run(tc.attrName, func(t *testing.T) {
			value, err := xattr.Get(filePath, tc.attrName)
			if err != nil {
				t.Fatalf("xattr.Get(%q) failed: %v", tc.attrName, err)
			}
			if string(value) != tc.expected {
				t.Errorf("Attribute value = %q; want %q", string(value), tc.expected)
			}
		})
	}
}
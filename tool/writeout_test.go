package tool

import (
	"strings"
	"testing"
)

func TestWriteOut(t *testing.T) {
	// Sample data to be used in the tests.
	sampleData := map[string]interface{}{
		"url_effective": "http://example.com",
		"http_code":     200,
		"time_total":    1.234567,
		"size_download": 1024,
	}

	testCases := []struct {
		name     string
		format   string
		data     map[string]interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "single string variable",
			format:   "%{url_effective}",
			data:     sampleData,
			expected: "http://example.com",
		},
		{
			name:     "single long variable with padding",
			format:   "%{http_code}",
			data:     sampleData,
			expected: "200",
		},
		{
			name:     "single time variable",
			format:   "%{time_total}",
			data:     sampleData,
			expected: "1.234567",
		},
		{
			name:     "mixed literals and variables",
			format:   "URL: %{url_effective} | RC: %{http_code}",
			data:     sampleData,
			expected: "URL: http://example.com | RC: 200",
		},
		{
			name:     "escaped percent sign",
			format:   "Downloaded %{size_download} bytes (100%%)",
			data:     sampleData,
			expected: "Downloaded 1024 bytes (100%)",
		},
		{
			name:     "unknown variable",
			format:   "This is %{nonsense} and should be ignored.",
			data:     sampleData,
			expected: "This is  and should be ignored.",
		},
		{
			name:     "variable with no data",
			format:   "Connect time: %{time_connect}",
			data:     sampleData, // time_connect is not in the map
			expected: "Connect time: 0.000000",
		},
		{
			name:     "no variables",
			format:   "Just a plain string.",
			data:     sampleData,
			expected: "Just a plain string.",
		},
		{
			name:    "invalid format string",
			format:  "Unmatched brace %{http_code",
			data:    sampleData,
			wantErr: true, // Expect an error for unmatched braces.
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var sb strings.Builder
			err := WriteOut(&sb, tc.format, tc.data)

			if (err != nil) != tc.wantErr {
				t.Fatalf("WriteOut() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				if result := sb.String(); result != tc.expected {
					t.Errorf("WriteOut() = %q; want %q", result, tc.expected)
				}
			}
		})
	}
}
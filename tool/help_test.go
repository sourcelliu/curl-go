package tool

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	testCases := []struct {
		name           string
		category       string
		expectToContain []string
	}{
		{
			name:     "default help",
			category: "",
			expectToContain: []string{
				"Usage: curl",
				"-o, --output <file>", // An important option
			},
		},
		{
			name:     "all help",
			category: "all",
			expectToContain: []string{
				"-d, --data <data>",
				"-H, --header <header>",
				"-I, --head",
				"-L, --location",
				"-o, --output <file>",
				"-u, --user <user:password>",
				"-v, --verbose",
			},
		},
		{
			name:     "http category",
			category: "http",
			expectToContain: []string{
				"HTTP and HTTPS protocol", // Category description
				"-H, --header <header>",   // HTTP-specific option
			},
		},
		{
			name:     "list categories",
			category: "category",
			expectToContain: []string{
				"auth",
				"connection",
				"curl",
				"http",
			},
		},
		{
			name:            "unknown category",
			category:        "nonsense",
			expectToContain: []string{"Unknown help category: nonsense"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintHelp(&buf, tc.category)
			output := buf.String()

			for _, expected := range tc.expectToContain {
				if !strings.Contains(output, expected) {
					t.Errorf("Help output for category %q did not contain %q.\nGot:\n%s",
						tc.category, expected, output)
				}
			}
		})
	}
}
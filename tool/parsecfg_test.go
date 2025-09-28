package tool

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []ConfigEntry
		wantErr  bool
	}{
		{
			name:     "simple option with equals",
			input:    `user-agent = "my-agent/1.0"`,
			expected: []ConfigEntry{{Option: "user-agent", Parameter: "my-agent/1.0"}},
		},
		{
			name:     "simple option with colon",
			input:    "output: output.html",
			expected: []ConfigEntry{{Option: "output", Parameter: "output.html"}},
		},
		{
			name:     "dashed option with space",
			input:    "--url http://example.com",
			expected: []ConfigEntry{{Option: "--url", Parameter: "http://example.com"}},
		},
		{
			name:     "short dashed option",
			input:    "-v",
			expected: []ConfigEntry{{Option: "-v", Parameter: ""}},
		},
		{
			name:     "comments and blank lines",
			input:    "\n# This is a comment\n\nverbose\n",
			expected: []ConfigEntry{{Option: "verbose", Parameter: ""}},
		},
		{
			name:     "quoted parameter with escapes",
			input:    `data = "hello\tworld\""`,
			expected: []ConfigEntry{{Option: "data", Parameter: "hello\tworld\""}},
		},
		{
			name:     "unquoted parameter with trailing comment",
			input:    "url = http://example.com # gets the site",
			expected: []ConfigEntry{{Option: "url", Parameter: "http://example.com"}},
		},
		{
			name: "multi-line config",
			input: `
# Set the user agent
user-agent = "Test Agent"

# And the URL
--url http://localhost/test

# A flag
--verbose

# Data with quotes and escapes
data = "a \"quoted\" string with a \\ backslash"
`,
			expected: []ConfigEntry{
				{Option: "user-agent", Parameter: "Test Agent"},
				{Option: "--url", Parameter: "http://localhost/test"},
				{Option: "--verbose", Parameter: ""},
				{Option: "data", Parameter: `a "quoted" string with a \ backslash`},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			result, err := ParseConfig(reader)

			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseConfig() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("ParseConfig() = %v; want %v", result, tc.expected)
			}
		})
	}
}
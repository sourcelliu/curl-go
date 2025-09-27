package tool

import "testing"

func TestBasename(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "unix path",
			path:     "/usr/local/bin/curl",
			expected: "curl",
		},
		{
			name:     "windows path",
			path:     "C:\\Users\\Jules\\curl.exe",
			expected: "curl.exe",
		},
		{
			name:     "mixed slashes",
			path:     "/home/user\\docs/file.txt",
			expected: "file.txt",
		},
		{
			name:     "no slashes",
			path:     "file.txt",
			expected: "file.txt",
		},
		{
			name:     "trailing forward slash",
			path:     "/usr/local/bin/",
			expected: "",
		},
		{
			name:     "trailing backslash",
			path:     "C:\\Users\\Jules\\",
			expected: "",
		},
		{
			name:     "empty path",
			path:     "",
			expected: "",
		},
		{
			name:     "root path",
			path:     "/",
			expected: "",
		},
		{
			name:     "multiple slashes",
			path:     "//a/b//c",
			expected: "c",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Basename(tc.path)
			if result != tc.expected {
				t.Errorf("Basename(%q) = %q; want %q", tc.path, result, tc.expected)
			}
		})
	}
}
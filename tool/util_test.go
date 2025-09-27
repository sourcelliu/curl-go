package tool

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTimeNow(t *testing.T) {
	// We can't test for an exact time, but we can check if it's recent.
	now := TimeNow()
	if time.Since(now) > time.Second {
		t.Errorf("TimeNow() returned a time that is more than 1 second in the past")
	}
}

func TestStricmp(t *testing.T) {
	testCases := []struct {
		name     string
		s1       string
		s2       string
		expected int
	}{
		{
			name:     "equal strings",
			s1:       "hello",
			s2:       "hello",
			expected: 0,
		},
		{
			name:     "case-insensitive equal",
			s1:       "Hello",
			s2:       "hello",
			expected: 0,
		},
		{
			name:     "s1 less than s2",
			s1:       "apple",
			s2:       "banana",
			expected: -1,
		},
		{
			name:     "s1 greater than s2",
			s1:       "cherry",
			s2:       "banana",
			expected: 1,
		},
		{
			name:     "empty strings",
			s1:       "",
			s2:       "",
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Stricmp(tc.s1, tc.s2)
			if result != tc.expected {
				t.Errorf("Stricmp(%q, %q) = %d; want %d", tc.s1, tc.s2, result, tc.expected)
			}
		})
	}
}

func TestTruncateFile(t *testing.T) {
	// Create a temporary file for the test.
	tempFile, err := os.CreateTemp(t.TempDir(), "truncate-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	// Write some data to the file.
	initialContent := "1234567890"
	if _, err := tempFile.WriteString(initialContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Truncate the file to 5 bytes.
	if err := TruncateFile(tempFile, 5); err != nil {
		t.Fatalf("TruncateFile() failed: %v", err)
	}

	// Check the new file size.
	info, err := tempFile.Stat()
	if err != nil {
		t.Fatalf("Failed to stat temp file: %v", err)
	}
	if info.Size() != 5 {
		t.Errorf("File size after truncate = %d; want 5", info.Size())
	}
}

func TestExecutableFile(t *testing.T) {
	// Get the path to the current running test executable.
	exePath, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() failed: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// Create a dummy file in the same directory.
	dummyFilename := "test-executable-file.txt"
	dummyPath := filepath.Join(exeDir, dummyFilename)
	if err := os.WriteFile(dummyPath, []byte("test"), 0600); err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}
	defer os.Remove(dummyPath) // Clean up the dummy file.

	t.Run("file exists", func(t *testing.T) {
		path, found := ExecutableFile(dummyFilename)
		if !found {
			t.Error("ExecutableFile() did not find an existing file")
		}
		if path != dummyPath {
			t.Errorf("ExecutableFile() path = %q; want %q", path, dummyPath)
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		_, found := ExecutableFile("non-existent-file.txt")
		if found {
			t.Error("ExecutableFile() found a non-existent file")
		}
	})
}
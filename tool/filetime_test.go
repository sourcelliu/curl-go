package tool

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileTime(t *testing.T) {
	// Create a temporary file for the test.
	tempFile, err := os.CreateTemp(t.TempDir(), "filetime-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFilename := tempFile.Name()
	tempFile.Close() // Close the file so it can be manipulated.

	// --- Test GetFileTime on a new file ---
	t.Run("get initial time", func(t *testing.T) {
		timestamp, err := GetFileTime(tempFilename)
		if err != nil {
			t.Fatalf("GetFileTime() failed on a new file: %v", err)
		}
		// The timestamp should be recent.
		if time.Since(time.Unix(timestamp, 0)) > 5*time.Second {
			t.Errorf("Initial file time is not recent")
		}
	})

	// --- Test SetFileTime and GetFileTime roundtrip ---
	t.Run("set and get time", func(t *testing.T) {
		// A specific, non-current timestamp. (Jan 1, 2020)
		expectedTimestamp := int64(1577836800)

		// Set the file time.
		if err := SetFileTime(expectedTimestamp, tempFilename); err != nil {
			t.Fatalf("SetFileTime() failed: %v", err)
		}

		// Get the file time back.
		actualTimestamp, err := GetFileTime(tempFilename)
		if err != nil {
			t.Fatalf("GetFileTime() failed after setting time: %v", err)
		}

		// Verify they match.
		if actualTimestamp != expectedTimestamp {
			t.Errorf("Timestamp mismatch: got %d, want %d", actualTimestamp, expectedTimestamp)
		}
	})

	// --- Test error handling for non-existent files ---
	t.Run("non-existent file", func(t *testing.T) {
		nonExistentFile := filepath.Join(t.TempDir(), "does-not-exist.txt")

		// Test GetFileTime error
		_, err := GetFileTime(nonExistentFile)
		if err == nil {
			t.Error("GetFileTime() did not return an error for a non-existent file")
		}

		// Test SetFileTime error
		err = SetFileTime(1, nonExistentFile)
		if err == nil {
			t.Error("SetFileTime() did not return an error for a non-existent file")
		}
	})
}
package tool

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestGetStderrWriter(t *testing.T) {
	t.Run("default stderr", func(t *testing.T) {
		writer, cleanup, err := GetStderrWriter("")
		defer cleanup()

		if err != nil {
			t.Fatalf("GetStderrWriter() with empty path returned an error: %v", err)
		}
		if writer != os.Stderr {
			t.Error("Expected os.Stderr, but got something else")
		}
	})

	t.Run("redirect to stdout", func(t *testing.T) {
		writer, cleanup, err := GetStderrWriter("-")
		defer cleanup()

		if err != nil {
			t.Fatalf("GetStderrWriter() with '-' returned an error: %v", err)
		}
		if writer != os.Stdout {
			t.Error("Expected os.Stdout, but got something else")
		}
	})

	t.Run("redirect to file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "stderr.log")

		writer, cleanup, err := GetStderrWriter(filePath)
		if err != nil {
			t.Fatalf("GetStderrWriter() with a file path returned an error: %v", err)
		}
		if writer == os.Stderr || writer == os.Stdout {
			t.Fatal("Expected a file writer, but got a standard stream")
		}
		defer cleanup() // This will close the file

		// Write to the file via the writer interface
		testString := "this is an error message"
		_, err = io.WriteString(writer, testString)
		if err != nil {
			t.Fatalf("Failed to write to the returned writer: %v", err)
		}

		// The cleanup function should close the file, so we can now read it.
		// Calling cleanup explicitly here to be clear, though defer would also work.
		cleanup()

		// Read the file content to verify
		content, readErr := os.ReadFile(filePath)
		if readErr != nil {
			t.Fatalf("Failed to read back the stderr file: %v", readErr)
		}
		if string(content) != testString {
			t.Errorf("File content mismatch: got %q, want %q", string(content), testString)
		}
	})

	t.Run("error on invalid path", func(t *testing.T) {
		// A path that cannot be created
		invalidPath := filepath.Join(t.TempDir(), "non-existent-dir", "stderr.log")

		_, _, err := GetStderrWriter(invalidPath)
		if err == nil {
			t.Error("GetStderrWriter() did not return an error for an invalid path")
		}
	})
}
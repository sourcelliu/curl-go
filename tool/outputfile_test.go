package tool

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateOutputFile(t *testing.T) {
	t.Run("ClobberAlways overwrites existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.txt")
		// Create an initial file with content.
		if err := os.WriteFile(filePath, []byte("initial"), 0644); err != nil {
			t.Fatalf("Failed to create initial file: %v", err)
		}

		file, _, err := CreateOutputFile(filePath, ClobberAlways, false)
		if err != nil {
			t.Fatalf("CreateOutputFile failed: %v", err)
		}
		defer file.Close()

		info, _ := file.Stat()
		if info.Size() != 0 {
			t.Error("File was not truncated when clobbering")
		}
	})

	t.Run("ClobberNever creates new file with suffix", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.txt")
		// Create the initial file.
		if err := os.WriteFile(filePath, []byte("initial"), 0644); err != nil {
			t.Fatalf("Failed to create initial file: %v", err)
		}

		_, finalName, err := CreateOutputFile(filePath, ClobberNever, false)
		if err != nil {
			t.Fatalf("CreateOutputFile failed: %v", err)
		}
		expectedName := filePath + ".1"
		if finalName != expectedName {
			t.Errorf("Expected new filename to be %q, but got %q", expectedName, finalName)
		}
	})

	t.Run("ClobberDefault with Content-Disposition fails if exists", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.txt")
		if err := os.WriteFile(filePath, []byte("initial"), 0644); err != nil {
			t.Fatalf("Failed to create initial file: %v", err)
		}

		_, _, err := CreateOutputFile(filePath, ClobberDefault, true)
		if err == nil {
			t.Error("Expected an error when file exists with Content-Disposition clobber mode, but got none")
		}
	})
}

func TestWriteCallback(t *testing.T) {
	t.Run("normal write", func(t *testing.T) {
		var buf bytes.Buffer
		config := &struct{ IsTTY, TerminalBinaryOK bool }{false, false}
		data := []byte("hello world")
		n, err := WriteCallback(&buf, data, config)

		if err != nil {
			t.Fatalf("WriteCallback failed: %v", err)
		}
		if n != len(data) {
			t.Errorf("Wrote %d bytes, want %d", n, len(data))
		}
		if buf.String() != "hello world" {
			t.Errorf("Buffer content is %q, want %q", buf.String(), "hello world")
		}
	})

	t.Run("binary to tty error", func(t *testing.T) {
		var buf bytes.Buffer
		config := &struct{ IsTTY, TerminalBinaryOK bool }{true, false}
		data := []byte("binary\x00data")
		_, err := WriteCallback(&buf, data, config)

		if err == nil {
			t.Error("Expected an error for binary output to TTY, but got none")
		}
	})

	t.Run("binary to tty allowed", func(t *testing.T) {
		var buf bytes.Buffer
		config := &struct{ IsTTY, TerminalBinaryOK bool }{true, true}
		data := []byte("binary\x00data")
		_, err := WriteCallback(&buf, data, config)

		if err != nil {
			t.Errorf("Did not expect an error when binary to TTY is allowed, but got: %v", err)
		}
	})
}
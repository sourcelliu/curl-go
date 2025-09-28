package tool

import (
	"os"
	"testing"
	"time"
)

func TestParseLong(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"positive", "12345", 12345, false},
		{"negative", "-54321", -54321, false},
		{"zero", "0", 0, false},
		{"invalid", "abc", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := ParseLong(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseLong(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}
			if !tc.wantErr && val != tc.expected {
				t.Errorf("ParseLong(%q) = %d, want %d", tc.input, val, tc.expected)
			}
		})
	}
}

func TestParseSecs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected time.Duration
		wantErr  bool
	}{
		{"integer", "30", 30 * time.Second, false},
		{"decimal", "1.5", 1500 * time.Millisecond, false},
		{"invalid", "xyz", 0, true},
		{"negative", "-10", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := ParseSecs(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseSecs(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}
			if !tc.wantErr && val != tc.expected {
				t.Errorf("ParseSecs(%q) = %v, want %v", tc.input, val, tc.expected)
			}
		})
	}
}

func TestFileToString(t *testing.T) {
	tempFile, err := os.CreateTemp(t.TempDir(), "test-file-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	content := "line 1\r\nline 2\nline 3"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	// Rewind the file to the beginning for reading
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp file: %v", err)
	}

	result, err := FileToString(tempFile)
	if err != nil {
		t.Fatalf("FileToString() failed: %v", err)
	}

	expected := "line 1line 2line 3"
	if result != expected {
		t.Errorf("FileToString() = %q; want %q", result, expected)
	}
}

func TestParseFTPFileMethod(t *testing.T) {
	testCases := []struct {
		input    string
		expected FTPFileMethod
		wantErr  bool
	}{
		{"multicwd", FTPMethodMultiCWD, false},
		{"NOCWD", FTPMethodNoCWD, false},
		{"singlecwd", FTPMethodSingleCWD, false},
		{"invalid", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			val, err := ParseFTPFileMethod(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseFTPFileMethod(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}
			if !tc.wantErr && val != tc.expected {
				t.Errorf("ParseFTPFileMethod(%q) = %v, want %v", tc.input, val, tc.expected)
			}
		})
	}
}

func TestParseDelegation(t *testing.T) {
	testCases := []struct {
		input    string
		expected GSSDelegation
		wantErr  bool
	}{
		{"none", DelegationNone, false},
		{"POLICY", DelegationPolicy, false},
		{"always", DelegationAlways, false},
		{"invalid", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			val, err := ParseDelegation(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseDelegation(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}
			if !tc.wantErr && val != tc.expected {
				t.Errorf("ParseDelegation(%q) = %v, want %v", tc.input, val, tc.expected)
			}
		})
	}
}
package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCurlInfoMain(t *testing.T) {
	// --- Setup: Build the binaries for testing ---
	tempDir := t.TempDir()

	// Build the default version (no tags)
	defaultBinaryPath := filepath.Join(tempDir, "curlinfo_default")
	buildCmdDefault := exec.Command("go", "build", "-o", defaultBinaryPath, ".")
	if err := buildCmdDefault.Run(); err != nil {
		t.Fatalf("Failed to build default binary: %v", err)
	}

	// Build the xattr version
	xattrBinaryPath := filepath.Join(tempDir, "curlinfo_xattr")
	buildCmdXattr := exec.Command("go", "build", "-tags", "xattr", "-o", xattrBinaryPath, ".")
	if err := buildCmdXattr.Run(); err != nil {
		t.Fatalf("Failed to build xattr binary: %v", err)
	}

	// --- Test Cases ---
	testCases := []struct {
		name          string
		binaryPath    string
		expectInOutput string
	}{
		{
			name:           "default build without xattr",
			binaryPath:     defaultBinaryPath,
			expectInOutput: "xattr                    : OFF",
		},
		{
			name:           "build with xattr tag",
			binaryPath:     xattrBinaryPath,
			expectInOutput: "xattr                    : ON",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(tc.binaryPath)
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Failed to run binary %s: %v", tc.binaryPath, err)
			}

			if !strings.Contains(string(output), tc.expectInOutput) {
				t.Errorf("Output did not contain expected string.\nWant to contain: %q\nGot:\n%s", tc.expectInOutput, string(output))
			}
		})
	}
}
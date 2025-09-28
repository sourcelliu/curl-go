package tool

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFindCurlRC(t *testing.T) {
	// --- Test Cases ---
	testCases := []struct {
		name         string
		setup        func(t *testing.T, dirs map[string]string) // Function to set up the test state
		expectedFile string                                     // The base name of the file we expect to find
		expectFound  bool
	}{
		{
			name: "found in CURL_HOME",
			setup: func(t *testing.T, dirs map[string]string) {
				// Create .curlrc in CURL_HOME
				path := filepath.Join(dirs["curlHome"], ".curlrc")
				if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create file in fakeCurlHome: %v", err)
				}
				// Also create one in home to ensure CURL_HOME has priority
				path = filepath.Join(dirs["home"], ".curlrc")
				if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create file in fakeHome: %v", err)
				}
			},
			expectedFile: ".curlrc",
			expectFound:  true,
		},
		{
			name: "found in user config dir",
			setup: func(t *testing.T, dirs map[string]string) {
				// No CURL_HOME, so it should check config dir
				path := filepath.Join(dirs["config"], "curl", ".curlrc")
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatalf("Failed to create nested config dir: %v", err)
				}
				if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create file in fakeConfig: %v", err)
				}
			},
			expectedFile: ".curlrc",
			expectFound:  true,
		},
		{
			name: "found in user home dir",
			setup: func(t *testing.T, dirs map[string]string) {
				// No CURL_HOME or config file, so it should check home dir
				path := filepath.Join(dirs["home"], ".curlrc")
				if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create file in fakeHome: %v", err)
				}
			},
			expectedFile: ".curlrc",
			expectFound:  true,
		},
		{
			name: "windows underscore fallback",
			setup: func(t *testing.T, dirs map[string]string) {
				if runtime.GOOS != "windows" {
					t.Skip("Skipping underscore test on non-Windows system")
				}
				path := filepath.Join(dirs["home"], "_curlrc")
				if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create _curlrc in fakeHome: %v", err)
				}
			},
			expectedFile: "_curlrc",
			expectFound:  true,
		},
		{
			name:         "not found anywhere",
			setup:        func(t *testing.T, dirs map[string]string) { /* No files created */ },
			expectedFile: "",
			expectFound:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- Setup an isolated temporary file structure for each test ---
			tempRoot := t.TempDir()
			dirs := map[string]string{
				"home":     filepath.Join(tempRoot, "home"),
				"config":   filepath.Join(tempRoot, "config"),
				"curlHome": filepath.Join(tempRoot, "curlhome"),
			}

			for _, dir := range dirs {
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("Failed to create fake directory: %v", err)
				}
			}

			// Set environment variables to control where the function looks.
			t.Setenv("CURL_HOME", "") // Unset by default
			if tc.name == "found in CURL_HOME" {
				t.Setenv("CURL_HOME", dirs["curlHome"])
			}
			t.Setenv("HOME", dirs["home"])
			t.Setenv("USERPROFILE", dirs["home"]) // For Windows home dir
			t.Setenv("XDG_CONFIG_HOME", dirs["config"])

			// Run the test-specific setup
			tc.setup(t, dirs)

			path, found := FindCurlRC()

			if found != tc.expectFound {
				t.Fatalf("FindCurlRC() found = %v; want %v", found, tc.expectFound)
			}

			if found {
				// Check just the filename, as the full path is temporary
				if filepath.Base(path) != tc.expectedFile {
					t.Errorf("FindCurlRC() found filename %q; want %q", filepath.Base(path), tc.expectedFile)
				}
			}
		})
	}
}
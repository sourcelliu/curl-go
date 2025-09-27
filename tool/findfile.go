package tool

import (
	"os"
	"path/filepath"
	"runtime"
)

// FindCurlRC searches for the curlrc file in standard locations and returns
// the path if found. This is a translation of the C function `findfile` from
// curl-src/src/tool_findfile.c, specialized for its primary use case.
//
// The search order is:
// 1. The path specified in the CURL_HOME environment variable.
// 2. The standard user configuration directory (e.g., ~/.config/curl/ on Linux).
// 3. The user's home directory.
//
// On Windows, it checks for both ".curlrc" and "_curlrc" in each location.
// On other systems, it only checks for ".curlrc".
func FindCurlRC() (string, bool) {
	// The C code has a complex list of environment variables to check.
	// Go's os.UserHomeDir and os.UserConfigDir provide cross-platform
	// ways to find the most important locations.

	// 1. Check CURL_HOME environment variable first.
	if curlHome := os.Getenv("CURL_HOME"); curlHome != "" {
		if path, found := checkFilesInDir(curlHome); found {
			return path, true
		}
	}

	// 2. Check the user's config directory.
	if configDir, err := os.UserConfigDir(); err == nil {
		// Look in a 'curl' subdirectory, as is common for XDG configs.
		curlConfigDir := filepath.Join(configDir, "curl")
		if path, found := checkFilesInDir(curlConfigDir); found {
			return path, true
		}
	}

	// 3. Check the user's home directory.
	if homeDir, err := os.UserHomeDir(); err == nil {
		if path, found := checkFilesInDir(homeDir); found {
			return path, true
		}
	}

	return "", false
}

// checkFilesInDir checks for ".curlrc" and, on Windows, "_curlrc" in the given directory.
func checkFilesInDir(dir string) (string, bool) {
	filenames := []string{".curlrc"}
	if runtime.GOOS == "windows" {
		filenames = append(filenames, "_curlrc")
	}

	for _, filename := range filenames {
		path := filepath.Join(dir, filename)
		if _, err := os.Stat(path); err == nil {
			// File exists.
			return path, true
		}
	}

	return "", false
}
package tool

import (
	"fmt"
	"io"
	"os"
)

// CleanupFunc is a function that can be deferred to clean up resources.
type CleanupFunc func()

// noOpCleanup does nothing. It's returned when no cleanup is necessary.
func noOpCleanup() {}

// GetStderrWriter determines the correct io.Writer for error messages based on a
// filename. It replaces the logic of `tool_set_stderr_file` from
// curl-src/src/tool_stderr.c.
//
// It returns an io.Writer and a cleanup function. The cleanup function
// MUST be called by the caller (e.g., with `defer`) to close any opened files.
//
// - If filename is empty, it returns os.Stderr.
// - If filename is "-", it returns os.Stdout.
// - Otherwise, it opens the specified file for writing.
func GetStderrWriter(filename string) (io.Writer, CleanupFunc, error) {
	if filename == "" {
		// Default case: use standard error.
		return os.Stderr, noOpCleanup, nil
	}

	if filename == "-" {
		// Special case: use standard output.
		return os.Stdout, noOpCleanup, nil
	}

	// Open the specified file for writing.
	// The C code uses "wt" which means write, text mode. Go's os.Create
	// truncates the file if it exists, which is the equivalent behavior.
	file, err := os.Create(filename)
	if err != nil {
		// The C code just prints a warning and continues using the old stderr.
		// In Go, it's more idiomatic to return the error and let the caller decide.
		return nil, nil, fmt.Errorf("failed to open stderr file %s: %w", filename, err)
	}

	// Return the file and a function to close it.
	cleanup := func() {
		file.Close()
	}
	return file, cleanup, nil
}
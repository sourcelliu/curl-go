package tool

import (
	"os"
	"path/filepath"
)

// CreateDirHierarchy creates the directory hierarchy for a given file path.
// This is a translation of the C function `create_dir_hierarchy` from
// curl-src/src/tool_dirhie.c.
//
// The C implementation manually iterates through the path and creates each
// directory. The idiomatic Go equivalent is to use `os.MkdirAll`, which
// handles this automatically and in a cross-platform way.
func CreateDirHierarchy(filePath string) error {
	// First, get the directory part of the file path.
	dir := filepath.Dir(filePath)

	// If dir is "." or empty, there's no directory to create.
	if dir == "" || dir == "." {
		return nil
	}

	// os.MkdirAll creates a directory path along with any necessary parents.
	// If the path is already a directory, MkdirAll does nothing and returns nil.
	// The permission 0755 is a standard, sensible default.
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		// The C code has a helper to print custom error messages based on errno.
		// In Go, it's more idiomatic to return the structured error itself,
		// allowing the caller to inspect it if needed.
		return err
	}

	return nil
}
package tool

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDirHierarchy(t *testing.T) {
	// Create a temporary directory to work in.
	tempDir := t.TempDir()

	t.Run("create nested directories", func(t *testing.T) {
		// Define a path for a file in a nested directory structure.
		nestedPath := filepath.Join(tempDir, "a", "b", "c", "file.txt")
		dirToCreate := filepath.Dir(nestedPath)

		// Call the function to create the hierarchy.
		if err := CreateDirHierarchy(nestedPath); err != nil {
			t.Fatalf("CreateDirHierarchy() failed: %v", err)
		}

		// Check if the directory was actually created.
		info, err := os.Stat(dirToCreate)
		if os.IsNotExist(err) {
			t.Errorf("Directory %q was not created", dirToCreate)
		} else if !info.IsDir() {
			t.Errorf("Path %q was created, but it is not a directory", dirToCreate)
		}
	})

	t.Run("path already exists", func(t *testing.T) {
		// Define a path and create it.
		existingPath := filepath.Join(tempDir, "exists", "file.txt")
		if err := CreateDirHierarchy(existingPath); err != nil {
			t.Fatalf("Initial CreateDirHierarchy() failed: %v", err)
		}

		// Call the function again on the same path. It should not fail.
		if err := CreateDirHierarchy(existingPath); err != nil {
			t.Errorf("CreateDirHierarchy() failed on an existing path: %v", err)
		}
	})

	t.Run("no directory to create", func(t *testing.T) {
		// A path with no directory component should do nothing and not error.
		simplePath := "file.txt"
		if err := CreateDirHierarchy(simplePath); err != nil {
			t.Errorf("CreateDirHierarchy() failed for a simple file path: %v", err)
		}
	})

	t.Run("file conflicts with directory path", func(t *testing.T) {
		// Create a file where a directory is supposed to go.
		conflictDir := filepath.Join(tempDir, "conflict-dir")
		if err := os.WriteFile(conflictDir, []byte("i am a file"), 0644); err != nil {
			t.Fatalf("Failed to create conflicting file: %v", err)
		}

		// Now try to create a hierarchy that uses that path as a directory.
		conflictingPath := filepath.Join(conflictDir, "file.txt")
		err := CreateDirHierarchy(conflictingPath)
		if err == nil {
			t.Error("CreateDirHierarchy() did not return an error when a file conflicted with a directory path")
		}
	})
}
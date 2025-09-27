package tool

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// ClobberMode defines the behavior when an output file already exists.
type ClobberMode int

const (
	ClobberDefault ClobberMode = iota // Default behavior (overwrite, unless from Content-Disposition)
	ClobberAlways                     // Always overwrite
	ClobberNever                      // Never overwrite, try numeric suffixes
)

// CreateOutputFile creates/opens a local file for writing, handling different
// overwrite modes. This is a translation of the C function
// `tool_create_output_file` from curl-src/src/tool_cb_wrt.c.
// It returns the opened file, the final filename used, and any error.
func CreateOutputFile(filename string, mode ClobberMode, isCdFilename bool) (*os.File, string, error) {
	// Determine if we should overwrite the file.
	shouldClobber := (mode == ClobberAlways) || (mode == ClobberDefault && !isCdFilename)

	if shouldClobber {
		file, err := os.Create(filename)
		return file, filename, err
	}

	// Logic for not overwriting (O_EXCL behavior)
	openFlags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	file, err := os.OpenFile(filename, openFlags, 0644)
	if err == nil {
		return file, filename, nil
	}

	// If the error is not "file exists", return it.
	if !os.IsExist(err) {
		return nil, "", err
	}

	// If mode is Never, try numeric suffixes.
	if mode == ClobberNever {
		for i := 1; i < 100; i++ {
			newFilename := fmt.Sprintf("%s.%d", filename, i)
			file, err = os.OpenFile(newFilename, openFlags, 0644)
			if err == nil {
				return file, newFilename, nil
			}
			if !os.IsExist(err) {
				return nil, "", err // A different error occurred
			}
		}
	}

	// If we're here, it's because the file exists and we shouldn't overwrite it.
	return nil, "", fmt.Errorf("file exists: %s", filename)
}

// WriteCallback is the Go equivalent of the `tool_write_cb` C function.
// It writes data to the provided writer and handles binary-to-tty detection.
// It returns the number of bytes written and any error.
func WriteCallback(writer io.Writer, data []byte, config *struct{ IsTTY, TerminalBinaryOK bool }) (int, error) {
	// Check for binary output to a terminal.
	if config != nil && config.IsTTY && !config.TerminalBinaryOK {
		if bytes.Contains(data, []byte{0}) {
			// In a real app, we'd use the Messager to print a warning.
			// For now, we return an error to signal the problem.
			return 0, fmt.Errorf("binary output can mess up your terminal")
		}
	}

	return writer.Write(data)
}
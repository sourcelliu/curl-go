package tool

import (
	"fmt"
	"io"
	"mime"
	"strings"
)

// HeaderProcessor handles processing of HTTP headers received during a transfer.
// It encapsulates the state and logic from `tool_cb_hdr.c`.
type HeaderProcessor struct {
	// --- Configuration ---
	// A list of writers to dump raw headers to (e.g., for --dump-header).
	HeaderWriters []io.Writer
	// A writer for the ETag (--etag-save).
	ETagWriter io.Writer
	// Whether to look for and extract a filename from Content-Disposition.
	HonorContentDisposition bool

	// --- State / Results ---
	// The filename extracted from a Content-Disposition header.
	FilenameFromDisposition string
}

// NewHeaderProcessor creates a new header processor.
func NewHeaderProcessor() *HeaderProcessor {
	return &HeaderProcessor{}
}

// Process handles a single header line. It is the Go equivalent of the C
// function `tool_header_cb`.
func (hp *HeaderProcessor) Process(headerLine string) error {
	// 1. Dump raw header if requested.
	if len(hp.HeaderWriters) > 0 {
		for _, w := range hp.HeaderWriters {
			if _, err := fmt.Fprint(w, headerLine); err != nil {
				return fmt.Errorf("failed to write header dump: %w", err)
			}
		}
	}

	// Trim whitespace for parsing.
	trimmedLine := strings.TrimSpace(headerLine)
	if trimmedLine == "" {
		return nil // End of headers
	}

	// 2. Check for ETag if requested.
	if hp.ETagWriter != nil && strings.HasPrefix(strings.ToLower(trimmedLine), "etag:") {
		etag := strings.TrimSpace(trimmedLine[5:])
		// The C code truncates the file first. We assume the writer is ready.
		if _, err := fmt.Fprintln(hp.ETagWriter, etag); err != nil {
			return fmt.Errorf("failed to write etag: %w", err)
		}
	}

	// 3. Check for Content-Disposition if requested.
	if hp.HonorContentDisposition && strings.HasPrefix(strings.ToLower(trimmedLine), "content-disposition:") {
		disposition := strings.TrimSpace(trimmedLine[20:])
		filename, err := parseContentDispositionFilename(disposition)
		if err == nil && filename != "" {
			hp.FilenameFromDisposition = filename
			hp.HonorContentDisposition = false // Found it, no need to look further.
		}
	}

	return nil
}

// parseContentDispositionFilename extracts the filename from a
// Content-Disposition header value. This is a much simpler and more robust
// implementation than the manual C parser `parse_filename`, thanks to Go's
// standard library.
func parseContentDispositionFilename(headerValue string) (string, error) {
	// The `mime.ParseMediaType` function is designed for Content-Type but works
	// perfectly for Content-Disposition as well.
	_, params, err := mime.ParseMediaType(headerValue)
	if err != nil {
		return "", err
	}
	// The filename is in the "filename" parameter.
	return params["filename"], nil
}
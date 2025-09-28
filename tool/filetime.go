package tool

import (
	"os"
	"time"
)

// GetFileTime retrieves the modification time of a file as a Unix timestamp.
// This is a translation of the C function `getfiletime` from
// curl-src/src/tool_filetime.c, lines 39-95.
//
// The C implementation uses platform-specific APIs (#ifdef _WIN32 ... #else ...).
// The idiomatic Go equivalent is to use `os.Stat`, which is cross-platform.
func GetFileTime(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	// The C function returns a curl_off_t (long) timestamp.
	// The Go equivalent is to get the modification time and convert it to
	// a Unix timestamp.
	return info.ModTime().Unix(), nil
}

// SetFileTime sets the modification time of a file from a Unix timestamp.
// This is a translation of the C function `setfiletime` from
// curl-src/src/tool_filetime.c, lines 98-158.
//
// The C implementation uses platform-specific APIs (utimes, utime, SetFileTime).
// The idiomatic Go equivalent is `os.Chtimes`, which is cross-platform.
func SetFileTime(timestamp int64, filename string) error {
	// The C function sets both access and modification time to the same value.
	mtime := time.Unix(timestamp, 0)
	atime := mtime // Keep access time the same as modification time.

	return os.Chtimes(filename, atime, mtime)
}
package tool

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TimeNow is a translation of the C function `tvrealnow` from
// curl-src/src/tool_util.c, lines 34-70.
//
// The C function returns the current time as a `struct timeval`, with
// different implementations for Windows and other systems. The idiomatic
// Go equivalent is `time.Now()`, which returns a `time.Time` object
// and is cross-platform.
func TimeNow() time.Time {
	return time.Now()
}

// Stricmp performs a case-insensitive string comparison, returning an integer
// (-1, 0, or 1). This is a translation of the C function `struplocompare`
// from curl-src/src/tool_util.c, lines 73-81.
//
// The C function also handles NULL pointers, which is not a direct concern
// for Go's non-nullable string type.
func Stricmp(s1, s2 string) int {
	return strings.Compare(strings.ToLower(s1), strings.ToLower(s2))
}

// The C file also contains `struplocompare4sort`, a wrapper for `qsort`.
// This is not needed in Go. To sort a slice of strings case-insensitively
// in Go, you can use the standard `sort` package like this:
//
//   sort.Slice(mySlice, func(i, j int) bool {
//     return Stricmp(mySlice[i], mySlice[j]) < 0
//   })

// TruncateFile truncates a file at a specific offset. This is a cross-platform
// Go equivalent of the Windows-specific C function `tool_ftruncate64` from
// curl-src/src/tool_util.c, lines 92-108.
func TruncateFile(file *os.File, offset int64) error {
	return file.Truncate(offset)
}

// ExecutableFile attempts to find a file with the given name in the same
// directory as the running executable. It returns the full path to the file
// and a boolean indicating if it exists. This is a cross-platform Go
// equivalent of the Windows-specific C function `tool_execpath` from
// curl-src/src/tool_util.c, lines 112-136.
func ExecutableFile(filename string) (string, bool) {
	exePath, err := os.Executable()
	if err != nil {
		return "", false
	}

	fullPath := filepath.Join(filepath.Dir(exePath), filename)

	if _, err := os.Stat(fullPath); err == nil {
		return fullPath, true
	}

	return "", false
}
package tool

import (
	"os"
	"strconv"

	"golang.org/x/term"
)

// GetTerminalColumns is a translation of the C function `get_terminal_columns`
// from curl-src/src/terminal.c.
//
// It returns the number of columns in the current terminal. It first checks
// the "COLUMNS" environment variable. If that is not set or invalid, it
// attempts to determine the width using a platform-specific system call.
// If all attempts fail, it returns a default width of 79.
func GetTerminalColumns() int {
	// Original C code logic from terminal.c, lines 45-53:
	//   char *colp = curl_getenv("COLUMNS");
	//   if(colp) { ... }
	if colp := os.Getenv("COLUMNS"); colp != "" {
		// The C code uses a custom string-to-number function.
		// strconv.Atoi is the standard Go equivalent.
		if width, err := strconv.Atoi(colp); err == nil {
			// The C code checks if the number is between 20 and 10000.
			// We'll replicate this check for consistency.
			if width > 20 && width < 10000 {
				return width
			}
		}
	}

	// Original C code logic from terminal.c, lines 56-89:
	// This section contains platform-specific ioctl/Windows API calls.
	// In Go, the idiomatic way to handle this is to use the golang.org/x/term
	// package, which provides a cross-platform way to get the terminal size.
	// We use os.Stdout.Fd() as the file descriptor for the terminal.
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		// The C code has a sanity check for the width.
		if width >= 0 && width < 10000 {
			return width
		}
	}

	// Original C code logic from terminal.c, line 92:
	//   if(!width)
	//     width = 79;
	//   return width;
	return 79 // Default width if all else fails.
}
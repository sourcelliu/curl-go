package main

import (
	"fmt"
	"runtime"
)

// This program is a Go translation of the intent of curl-src/src/curlinfo.c.
// The C program uses the preprocessor to report on features enabled at
// compile time. Go doesn't have a preprocessor, so we use a combination of
// the `runtime` package and build tags to achieve a similar result.

func main() {
	// A map to hold features and their status.
	// The status is determined by Go's standard library capabilities and
	// build tags (like `xattrEnabled`).
	features := map[string]bool{
		"cookies":               true, // Go's net/http/cookiejar
		"DoH":                   true, // Go's net.Resolver can be configured for DoH
		"HTTP-auth":             true, // Handled by net/http
		"Mime":                  true, // Go's mime/multipart
		"netrc":                 true, // Can be implemented with netrc package
		"proxy":                 true, // Handled by net/http.Transport
		"shuffle-dns":           true, // Go's net.Resolver shuffles by default
		"large-time":            true, // Go's time.Time is 64-bit
		"large-size":            true, // Go's int/uint are 64-bit on 64-bit systems
		"xattr":                 xattrEnabled,
		"win32-ca-searchpath":   runtime.GOOS == "windows",
	}

	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("\nFeatures:")

	// Print the status of each feature.
	for name, enabled := range features {
		status := "OFF"
		if enabled {
			status = "ON"
		}
		fmt.Printf("%-25s: %s\n", name, status)
	}
}
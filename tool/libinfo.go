package tool

import (
	"runtime"
	"sync"
)

// Info holds information about the build-time and runtime capabilities
// of the translated curl-go application.
// This struct and the GetInfo function serve as the Go equivalent for the C
// file `tool_libinfo.c`, which queries the linked libcurl library.
type Info struct {
	Version   string
	GoVersion string
	OS        string
	Arch      string
	Protocols []string
	Features  map[string]bool
}

var (
	infoOnce sync.Once
	curlInfo *Info
)

// GetInfo returns a singleton instance of the Info struct, containing
// information about the application's capabilities. It mimics the C function
// `get_libcurl_info` by populating the info on the first call.
func GetInfo() *Info {
	infoOnce.Do(func() {
		// In a real application, this would be the version of our tool.
		const toolVersion = "0.1.0-alpha"

		// The Go net/http client supports these protocols by default.
		protocols := []string{
			"http",
			"httpss",
		}

		// Map C feature names to their Go equivalents.
		// The status is determined by standard library features or build tags.
		features := map[string]bool{
			"SSL":         true, // Go's crypto/tls is always available.
			"HTTP2":       true, // Go's net/http client supports HTTP/2 by default.
			"IPv6":        true, // Go's net package supports IPv6.
			"libz":        true, // Go's compress/gzip and compress/zlib.
			"brotli":      false, // Would require a third-party library.
			"zstd":        false, // Would require a third-party library.
			"Unicode":     true, // Go strings are UTF-8 by default.
			"threadsafe":  true, // Go has built-in concurrency.
			"xattr":       XattrEnabled, // This comes from the build tags.
		}

		curlInfo = &Info{
			Version:   toolVersion,
			GoVersion: runtime.Version(),
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
			Protocols: protocols,
			Features:  features,
		}
	})
	return curlInfo
}
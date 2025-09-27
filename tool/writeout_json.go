package tool

import (
	"encoding/json"
	"io"
	"net/http"
	"runtime"
)

// WriteOutJSON formats the provided data map as a single JSON object and writes
// it to the writer. This is a translation of the C function `ourWriteOutJSON`
// from curl-src/src/tool_writeout_json.c.
//
// The C implementation manually constructs the JSON string, which is complex
// and error-prone. The idiomatic Go equivalent is to use the standard
// `encoding/json` package to marshal a map directly.
func WriteOutJSON(writer io.Writer, data map[string]interface{}) error {
	// Create a copy of the map to avoid modifying the caller's data.
	outData := make(map[string]interface{}, len(data)+1)
	for k, v := range data {
		outData[k] = v
	}

	// Add the curl version information, similar to the C implementation.
	// This was a special case added at the end in the C code.
	outData["curl_version"] = runtime.Version()

	// Use an encoder for efficient writing to the stream.
	encoder := json.NewEncoder(writer)
	// encoder.SetIndent("", "  ") // Optionally pretty-print the JSON

	return encoder.Encode(outData)
}

// HeaderJSON formats HTTP headers into a JSON object. This is a translation of
// the C function `headerJSON` from curl-src/src/tool_writeout_json.c.
//
// The C implementation iterates through headers and manually builds the JSON.
// The idiomatic Go equivalent is to simply marshal the `http.Header` map.
func HeaderJSON(writer io.Writer, headers http.Header) error {
	encoder := json.NewEncoder(writer)
	// encoder.SetIndent("", "  ") // Optionally pretty-print the JSON

	return encoder.Encode(headers)
}

// The C file also contains `jsonquoted` and `jsonWriteString` for manual
// string escaping. These are not needed in Go, as the `encoding/json`
// package handles all necessary string escaping automatically and safely.
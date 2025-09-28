package tool

// Strdup is a translation of the C function `strdup` from
// curl-src/src/tool_strdup.c.
//
// In C, strdup allocates new memory for a copy of a string. In Go, strings
// are immutable and managed by the runtime, so a direct copy is not needed.
// This function is provided for completeness of the translation.
//
// Original C code from tool_strdup.c, lines 26-42.
func Strdup(s string) string {
	// A simple string assignment in Go is sufficient and efficient. It copies
	// the string header (pointer and length), not the underlying byte array.
	// The original string's data is immutable and safe to share.
	return s
}

// Memdup0 is a translation of the C function `memdup0` from
// curl-src/src/tool_strdup.c.
//
// The C function copies a memory buffer of a given length and ensures it is
// null-terminated. The idiomatic Go equivalent is to convert a byte slice
// to a string, which creates a new, immutable string by copying the bytes.
//
// Original C code from tool_strdup.c, lines 44-53.
func Memdup0(data []byte) string {
	// string(data) creates a new string by copying the bytes from the slice.
	return string(data)
}
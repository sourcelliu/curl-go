//go:build !xattr

package tool

// WriteXattr is a dummy implementation for when the 'xattr' build tag is not used.
// It does nothing and returns nil.
func WriteXattr(filePath string, data map[string]interface{}) error {
	// This is a no-op because extended attribute support was not enabled at build time.
	return nil
}
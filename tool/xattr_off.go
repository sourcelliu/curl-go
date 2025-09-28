//go:build !xattr

package tool

// XattrEnabled is false when the 'xattr' build tag is not used.
const XattrEnabled = false
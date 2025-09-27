//go:build !xattr

package main

// xattrEnabled is false when the 'xattr' build tag is not used.
const xattrEnabled = false
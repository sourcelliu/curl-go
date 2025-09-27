//go:build xattr

package tool

import (
	"fmt"
	"net/url"

	"github.com/pkg/xattr"
)

// stripCredentials removes user credentials from a URL string.
// This is a translation of the C function `stripcredentials` from
// curl-src/src/tool_xattr.c.
func stripCredentials(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	u.User = nil // Remove username and password
	return u.String(), nil
}

// WriteXattr sets extended attributes on a file from a map of metadata.
// This is a translation of the C function `fwrite_xattr` from
// curl-src/src/tool_xattr.c.
//
// It takes a file path and a map of data. The keys in the map correspond
// to the curl write-out variable names (e.g., "content_type").
func WriteXattr(filePath string, data map[string]interface{}) error {
	// Mapping from curl variable names to xattr attribute names.
	// From https://freedesktop.org/wiki/CommonExtendedAttributes/
	mappings := map[string]string{
		"content_type": "user.mime_type",
		"referer":      "user.xdg.referrer.url",
	}

	// Set the creator attribute, as curl does.
	if err := xattr.Set(filePath, "user.creator", []byte("curl-translation-go")); err != nil {
		return fmt.Errorf("failed to set creator xattr: %w", err)
	}

	// Set attributes based on the provided data.
	for varName, attrName := range mappings {
		if value, ok := data[varName].(string); ok && value != "" {
			if err := xattr.Set(filePath, attrName, []byte(value)); err != nil {
				return fmt.Errorf("failed to set xattr %s: %w", attrName, err)
			}
		}
	}

	// Set the origin URL after stripping credentials.
	if rawURL, ok := data["url_effective"].(string); ok && rawURL != "" {
		cleanURL, err := stripCredentials(rawURL)
		if err != nil {
			return fmt.Errorf("failed to strip credentials from URL: %w", err)
		}
		if err := xattr.Set(filePath, "user.xdg.origin.url", []byte(cleanURL)); err != nil {
			return fmt.Errorf("failed to set origin url xattr: %w", err)
		}
	}

	return nil
}
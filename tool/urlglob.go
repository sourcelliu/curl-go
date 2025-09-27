package tool

import (
	"fmt"
	"strconv"
	"strings"
)

// globPattern represents a single component of a URL glob, which can be a
// literal string, a set of options, or a character/numeric range.
type globPattern struct {
	options []string
}

// ExpandURLGlob takes a URL with globbing patterns and returns a slice of all
// expanded URLs. This is the primary entry point for the globbing functionality.
// It is a translation of the C function `glob_url` and its helpers from
// curl-src/src/tool_urlglob.c.
func ExpandURLGlob(pattern string) ([]string, error) {
	patterns, err := parse(pattern)
	if err != nil {
		return nil, err
	}

	if len(patterns) == 0 {
		return []string{""}, nil
	}

	var results []string
	generate(&results, patterns, "", 0)
	return results, nil
}

// parse iterates through the URL pattern and breaks it down into a series of
// globPattern structs.
func parse(pattern string) ([]globPattern, error) {
	var patterns []globPattern
	var currentLiteral strings.Builder

	for i := 0; i < len(pattern); i++ {
		char := pattern[i]

		switch char {
		case '\\':
			// Handle escaped characters
			i++
			if i < len(pattern) {
				currentLiteral.WriteByte(pattern[i])
			}
		case '{':
			// Flush the current literal part
			if currentLiteral.Len() > 0 {
				patterns = append(patterns, globPattern{options: []string{currentLiteral.String()}})
				currentLiteral.Reset()
			}
			// Parse the set
			set, remainder, err := parseSet(pattern[i+1:])
			if err != nil {
				return nil, err
			}
			patterns = append(patterns, set)
			pattern = remainder
			i = -1 // Restart loop at the beginning of the new remainder
		case '[':
			// Flush the current literal part
			if currentLiteral.Len() > 0 {
				patterns = append(patterns, globPattern{options: []string{currentLiteral.String()}})
				currentLiteral.Reset()
			}
			// Parse the range
			rangePattern, remainder, err := parseRange(pattern[i+1:])
			if err != nil {
				return nil, err
			}
			patterns = append(patterns, rangePattern)
			pattern = remainder
			i = -1 // Restart loop
		default:
			currentLiteral.WriteByte(char)
		}
	}

	// Add any remaining literal part
	if currentLiteral.Len() > 0 {
		patterns = append(patterns, globPattern{options: []string{currentLiteral.String()}})
	}

	return patterns, nil
}

// parseSet handles patterns like {a,b,c}.
func parseSet(pattern string) (globPattern, string, error) {
	end := strings.IndexByte(pattern, '}')
	if end == -1 {
		return globPattern{}, "", fmt.Errorf("unmatched brace in glob")
	}
	content := pattern[:end]
	options := strings.Split(content, ",")
	return globPattern{options: options}, pattern[end+1:], nil
}

// parseRange handles patterns like [a-z] or [0-9].
func parseRange(pattern string) (globPattern, string, error) {
	end := strings.IndexByte(pattern, ']')
	if end == -1 {
		return globPattern{}, "", fmt.Errorf("unmatched bracket in glob")
	}
	content := pattern[:end]
	parts := strings.Split(content, "-")
	if len(parts) != 2 {
		return globPattern{}, "", fmt.Errorf("invalid range format: %s", content)
	}

	start, endRange := parts[0], parts[1]
	var options []string

	// Check if it's a numeric range
	if startNum, err1 := strconv.Atoi(start); err1 == nil {
		if endNum, err2 := strconv.Atoi(endRange); err2 == nil {
			if startNum > endNum {
				return globPattern{}, "", fmt.Errorf("invalid numeric range: start > end")
			}
			for i := startNum; i <= endNum; i++ {
				options = append(options, strconv.Itoa(i))
			}
			return globPattern{options: options}, pattern[end+1:], nil
		}
	}

	// Assume character range
	if len(start) != 1 || len(endRange) != 1 {
		return globPattern{}, "", fmt.Errorf("invalid character range format")
	}
	startChar, endChar := start[0], endRange[0]
	if startChar > endChar {
		return globPattern{}, "", fmt.Errorf("invalid character range: start > end")
	}
	for c := startChar; c <= endChar; c++ {
		options = append(options, string(c))
	}
	return globPattern{options: options}, pattern[end+1:], nil
}

// generate recursively builds the final URL strings from the parsed patterns.
func generate(results *[]string, patterns []globPattern, current string, index int) {
	if index == len(patterns) {
		*results = append(*results, current)
		return
	}

	for _, option := range patterns[index].options {
		generate(results, patterns, current+option, index+1)
	}
}
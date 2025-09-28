package tool

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// ConfigEntry represents a single parsed option and its parameter from a config file.
type ConfigEntry struct {
	Option    string
	Parameter string
}

// unslashQuote is a helper function that translates the C function `unslashquote`
// from curl-src/src/tool_parsecfg.c, lines 74-101.
// It parses a string, handling backslash-escaped characters, and stops at the
// first non-escaped double quote. It returns the unquoted string and the
// remainder of the input string.
func unslashQuote(line string) (string, string) {
	var sb strings.Builder
	var i int
	for i = 0; i < len(line); i++ {
		char := line[i]
		if char == '\\' {
			i++ // Move to the character after the backslash
			if i < len(line) {
				escapedChar := line[i]
				switch escapedChar {
				case 't':
					sb.WriteByte('\t')
				case 'n':
					sb.WriteByte('\n')
				case 'r':
					sb.WriteByte('\r')
				case 'v':
					sb.WriteByte('\v')
				default:
					sb.WriteByte(escapedChar)
				}
			}
		} else if char == '"' {
			// End of quoted string
			return sb.String(), line[i+1:]
		} else {
			sb.WriteByte(char)
		}
	}
	// Reached end of line without a closing quote
	return sb.String(), ""
}

// ParseConfig reads from an io.Reader and parses it as a curl config file.
// It is a translation of the C function `parseconfig` from
// curl-src/src/tool_parsecfg.c, lines 105-274.
// It returns a slice of ConfigEntry structs or an error.
func ParseConfig(reader io.Reader) ([]ConfigEntry, error) {
	var entries []ConfigEntry
	scanner := bufio.NewScanner(reader)
	lineno := 0

	for scanner.Scan() {
		lineno++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		option := line
		var parameter string

		dashedOption := strings.HasPrefix(option, "-")

		// Find the end of the option keyword
		splitIndex := strings.IndexFunc(line, func(r rune) bool {
			if unicode.IsSpace(r) {
				return true
			}
			// Separators are only valid if the option is not dashed
			if !dashedOption && (r == '=' || r == ':') {
				return true
			}
			return false
		})

		if splitIndex != -1 {
			option = line[:splitIndex]
			paramPart := line[splitIndex:]

			// Find the start of the parameter, skipping whitespace and separators
			paramStartIndex := strings.IndexFunc(paramPart, func(r rune) bool {
				return !unicode.IsSpace(r) && (dashedOption || (r != '=' && r != ':'))
			})

			if paramStartIndex != -1 {
				paramPart = paramPart[paramStartIndex:]
				if strings.HasPrefix(paramPart, "\"") {
					// Quoted parameter
					parameter, _ = unslashQuote(paramPart[1:])
				} else {
					// Unquoted parameter is the first word
					endParamIndex := strings.IndexFunc(paramPart, unicode.IsSpace)
					if endParamIndex != -1 {
						parameter = paramPart[:endParamIndex]
						// The C code warns about unquoted whitespace here.
					} else {
						parameter = paramPart
					}
				}
			}
		}
		entries = append(entries, ConfigEntry{Option: option, Parameter: parameter})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	return entries, nil
}
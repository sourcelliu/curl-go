package tool

import (
	"fmt"
	"strings"
)

// FormPartType defines the type of a form part's content.
type FormPartType int

const (
	PartTypeLiteral  FormPartType = iota // e.g., name=value
	PartTypeFile                       // e.g., name=@file.txt
	PartTypeDataFile                   // e.g., name=<file.txt
)

// FormPart represents a single part of a multipart form submission.
// This struct is the Go equivalent of the C `tool_mime` struct and its data.
type FormPart struct {
	Name        string
	Value       string // Can be literal data or a filename
	Type        FormPartType
	ContentType string
	Filename    string // Can be used to override the original filename
	Encoder     string
	Headers     []string
}

// parser holds the state for parsing a form string. It allows us to process
// the string token by token, similar to how the C code uses pointers.
type parser struct {
	input string
	pos   int
}

func (p *parser) skipSpaces() {
	for p.pos < len(p.input) && (p.input[p.pos] == ' ' || p.input[p.pos] == '\t') {
		p.pos++
	}
}

// getWord extracts a word, handling quotes. The word ends at a separator in `endChars`.
func (p *parser) getWord(endChars string) string {
	p.skipSpaces()
	start := p.pos

	if p.pos < len(p.input) && p.input[p.pos] == '"' {
		p.pos++ // skip leading quote
		var sb strings.Builder
		for p.pos < len(p.input) {
			char := p.input[p.pos]
			if char == '\\' && p.pos+1 < len(p.input) {
				p.pos++
				sb.WriteByte(p.input[p.pos])
			} else if char == '"' {
				p.pos++ // skip trailing quote
				return sb.String()
			} else {
				sb.WriteByte(char)
			}
			p.pos++
		}
		p.pos = start // Unclosed quote, treat as literal from the start.
	}

	end := len(p.input)
	for i := p.pos; i < len(p.input); i++ {
		if strings.ContainsRune(endChars, rune(p.input[i])) {
			end = i
			break
		}
	}
	word := p.input[p.pos:end]
	p.pos = end
	return strings.TrimSpace(word)
}

// ParseFormString parses a curl-style form string (e.g., "name=value;type=...").
// It returns a slice of FormPart structs, as one string can define multiple parts.
// This is a translation of the C function `formparse`.
// Note: The C version also handles multipart boundaries with '(' and ')', which
// is omitted here for simplicity as it relates to higher-level state management.
func ParseFormString(input string) ([]*FormPart, error) {
	name, content, found := strings.Cut(input, "=")
	if !found {
		return nil, fmt.Errorf("invalid form string: missing '='")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("form part has no name")
	}

	var parts []*FormPart
	p := &parser{input: content}

	for {
		part := &FormPart{Name: name}
		p.skipSpaces()
		if p.pos >= len(p.input) {
			if len(parts) > 0 {
				break
			} // Trailing comma case
			return nil, fmt.Errorf("form part has no value")
		}

		// Determine value type
		if p.input[p.pos] == '@' {
			part.Type = PartTypeFile
			p.pos++
		} else if p.input[p.pos] == '<' {
			part.Type = PartTypeDataFile
			p.pos++
		} else {
			part.Type = PartTypeLiteral
		}

		part.Value = p.getWord(";,")
		part.Filename = part.Value // Default filename is the value itself

		// Parse semicolon-separated attributes
		for p.pos < len(p.input) && p.input[p.pos] == ';' {
			p.pos++ // skip semicolon
			attr := p.getWord("=")
			if p.pos < len(p.input) && p.input[p.pos] == '=' {
				p.pos++ // skip equals
				val := p.getWord(";,")
				switch strings.ToLower(attr) {
				case "type":
					part.ContentType = val
				case "filename":
					part.Filename = val
				case "encoder":
					part.Encoder = val
				case "headers":
					// Simplified: doesn't support reading headers from a file (@)
					part.Headers = append(part.Headers, val)
				}
			}
		}
		parts = append(parts, part)

		if p.pos >= len(p.input) || p.input[p.pos] != ',' {
			break // No more parts
		}
		p.pos++ // Skip comma

		if part.Type == PartTypeLiteral {
			return nil, fmt.Errorf("literal form parts cannot be comma-separated")
		}
	}

	return parts, nil
}
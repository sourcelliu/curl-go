package tool

import (
	"fmt"
	"io"
	"strings"
)

// VariableType defines the type of a --write-out variable.
type VariableType int

const (
	VarTypeString VariableType = iota
	VarTypeLong
	VarTypeOffset
	VarTypeTime
	VarTypeJSON
	VarTypeSpecial
)

// WriteOutVariable holds metadata about a single --write-out variable.
// Based on the C struct `writeoutvar`.
type WriteOutVariable struct {
	Name string
	Type VariableType
}

// variables is a map of all supported --write-out variables.
// This is a translation of the static C array `variables` from
// curl-src/src/tool_writeout.c. A map is used for efficient lookups.
var variables = map[string]WriteOutVariable{
	"content_type":          {Name: "content_type", Type: VarTypeString},
	"filename_effective":    {Name: "filename_effective", Type: VarTypeString},
	"ftp_entry_path":        {Name: "ftp_entry_path", Type: VarTypeString},
	"http_code":             {Name: "http_code", Type: VarTypeLong},
	"http_connect":          {Name: "http_connect", Type: VarTypeLong},
	"http_version":          {Name: "http_version", Type: VarTypeString},
	"local_ip":              {Name: "local_ip", Type: VarTypeString},
	"local_port":            {Name: "local_port", Type: VarTypeLong},
	"method":                {Name: "method", Type: VarTypeString},
	"num_connects":          {Name: "num_connects", Type: VarTypeLong},
	"num_redirects":         {Name: "num_redirects", Type: VarTypeLong},
	"redirect_url":          {Name: "redirect_url", Type: VarTypeString},
	"referer":               {Name: "referer", Type: VarTypeString},
	"remote_ip":             {Name: "remote_ip", Type: VarTypeString},
	"remote_port":           {Name: "remote_port", Type: VarTypeLong},
	"response_code":         {Name: "response_code", Type: VarTypeLong},
	"scheme":                {Name: "scheme", Type: VarTypeString},
	"size_download":         {Name: "size_download", Type: VarTypeOffset},
	"size_header":           {Name: "size_header", Type: VarTypeLong},
	"size_request":          {Name: "size_request", Type: VarTypeLong},
	"size_upload":           {Name: "size_upload", Type: VarTypeOffset},
	"speed_download":        {Name: "speed_download", Type: VarTypeOffset},
	"speed_upload":          {Name: "speed_upload", Type: VarTypeOffset},
	"ssl_verify_result":     {Name: "ssl_verify_result", Type: VarTypeLong},
	"time_appconnect":       {Name: "time_appconnect", Type: VarTypeTime},
	"time_connect":          {Name: "time_connect", Type: VarTypeTime},
	"time_namelookup":       {Name: "time_namelookup", Type: VarTypeTime},
	"time_posttransfer":     {Name: "time_posttransfer", Type: VarTypeTime},
	"time_pretransfer":      {Name: "time_pretransfer", Type: VarTypeTime},
	"time_redirect":         {Name: "time_redirect", Type: VarTypeTime},
	"time_starttransfer":    {Name: "time_starttransfer", Type: VarTypeTime},
	"time_total":            {Name: "time_total", Type: VarTypeTime},
	"url_effective":         {Name: "url_effective", Type: VarTypeString},
	// Special variables that need custom handling
	"json":        {Name: "json", Type: VarTypeJSON},
	"stderr":      {Name: "stderr", Type: VarTypeSpecial},
	"stdout":      {Name: "stdout", Type: VarTypeSpecial},
}

// WriteOut parses a format string and substitutes variables from the data map.
// This is a translation of the C function `ourWriteOut` from
// curl-src/src/tool_writeout.c.
func WriteOut(writer io.Writer, format string, data map[string]interface{}) error {
	var i int
	for i < len(format) {
		char := format[i]
		if char == '%' && i+1 < len(format) {
			i++ // Move past the '%'
			if format[i] == '%' {
				// Escaped '%%'
				fmt.Fprint(writer, "%")
				i++
			} else if format[i] == '{' {
				// Variable substitution %{...}
				i++ // Move past the '{'
				end := strings.IndexByte(format[i:], '}')
				if end == -1 {
					return fmt.Errorf("unmatched brace in write-out format")
				}
				varName := format[i : i+end]
				i += end + 1 // Move past the '}'

				if v, ok := variables[varName]; ok {
					val, dataOk := data[v.Name]
					if !dataOk {
						// In curl, this often prints 0 or an empty string.
						// We'll print a default value based on type.
						switch v.Type {
						case VarTypeLong, VarTypeOffset:
							fmt.Fprint(writer, "0")
						case VarTypeTime:
							fmt.Fprint(writer, "0.000000")
						}
						continue
					}
					// Format the value based on its type
					switch v.Type {
					case VarTypeString:
						fmt.Fprintf(writer, "%s", val)
					case VarTypeLong:
						// Special case for http_code
						if v.Name == "http_code" {
							fmt.Fprintf(writer, "%03d", val)
						} else {
							fmt.Fprintf(writer, "%d", val)
						}
					case VarTypeOffset:
						fmt.Fprintf(writer, "%d", val)
					case VarTypeTime:
						if f, ok := val.(float64); ok {
							fmt.Fprintf(writer, "%.6f", f)
						} else {
							fmt.Fprint(writer, "0.000000")
						}
					// TODO: Implement special cases like json, stderr, etc.
					}
				} else {
					// Unknown variable, curl prints a warning to stderr
					// but we'll just ignore it for now.
				}
			} else {
				// Invalid format, print literally
				fmt.Fprintf(writer, "%%%c", format[i])
				i++
			}
		} else {
			// Regular character
			fmt.Fprint(writer, string(char))
			i++
		}
	}
	return nil
}
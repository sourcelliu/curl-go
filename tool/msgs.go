package tool

import (
	"fmt"
	"io"
	"strings"
)

// Messager handles printing formatted, word-wrapped messages to a writer.
// It is a translation of the logic in curl-src/src/tool_msgs.c.
// The C implementation uses global state (the 'global' struct), which is
// encapsulated here for better testability and modularity.
type Messager struct {
	// Writer is the output stream, typically os.Stderr.
	Writer io.Writer
	// Silent corresponds to the --silent flag.
	Silent bool
	// Trace corresponds to the --verbose flag (global->tracetype).
	Trace bool
	// ShowError corresponds to the --show-error flag.
	ShowError bool
}

// NewMessager creates a new Messager with a given writer and configuration.
func NewMessager(writer io.Writer, silent, trace, showError bool) *Messager {
	return &Messager{
		Writer:    writer,
		Silent:    silent,
		Trace:     trace,
		ShowError: showError,
	}
}

// voutf is a private helper that formats and prints a message with a given
// prefix, handling word wrapping. It is a translation of the static C
// function `voutf` from curl-src/src/tool_msgs.c, lines 42-93.
func (m *Messager) voutf(prefix string, format string, args ...interface{}) {
	// Get the full, formatted message.
	message := fmt.Sprintf(format, args...)

	// Get terminal width and calculate available space.
	terminalWidth := GetTerminalColumns()
	availableWidth := terminalWidth - len(prefix)
	if availableWidth < 1 {
		availableWidth = 1 // Avoid negative or zero width.
	}

	// Split the message by newlines to handle multi-line inputs gracefully.
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		// Process each line, wrapping it as needed.
		for len(line) > 0 {
			fmt.Fprint(m.Writer, prefix)

			var cut int
			if len(line) > availableWidth {
				// Find the last space to break the line cleanly.
				lastSpace := strings.LastIndex(line[:availableWidth+1], " ")
				if lastSpace > 0 {
					cut = lastSpace
				} else {
					// No space found, so we have to cut mid-word.
					cut = availableWidth
				}
				fmt.Fprintln(m.Writer, line[:cut])
				line = strings.TrimLeft(line[cut:], " ")
			} else {
				fmt.Fprintln(m.Writer, line)
				line = ""
			}
		}
	}
}

// Notef prints a 'note' message if tracing (verbose) is enabled.
// Translation of `notef` from curl-src/src/tool_msgs.c, lines 100-108.
func (m *Messager) Notef(format string, args ...interface{}) {
	if m.Trace && !m.Silent {
		m.voutf("Note: ", format, args...)
	}
}

// Warnf prints a 'warning' message unless silenced.
// Translation of `warnf` from curl-src/src/tool_msgs.c, lines 111-119.
func (m *Messager) Warnf(format string, args ...interface{}) {
	if !m.Silent {
		m.voutf("Warning: ", format, args...)
	}
}

// Helpf prints a help-related error message.
// This is a simplified translation of `helpf` from curl-src/src/tool_msgs.c,
// lines 122-138.
func (m *Messager) Helpf(format string, args ...interface{}) {
	if format != "" {
		message := fmt.Sprintf(format, args...)
		fmt.Fprintf(m.Writer, "curl: %s\n", message)
	}
	fmt.Fprintln(m.Writer, "curl: try 'curl --help' or 'curl --manual' for more information")
}

// Errorf prints an error message unless silenced. It can be forced with
// the --show-error flag.
// Translation of `errorf` from curl-src/src/tool_msgs.c, lines 141-150.
func (m *Messager) Errorf(format string, args ...interface{}) {
	if !m.Silent || m.ShowError {
		m.voutf("curl: ", format, args...)
	}
}
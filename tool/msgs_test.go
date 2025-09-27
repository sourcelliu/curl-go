package tool

import (
	"bytes"
	"strings"
	"testing"
)

func TestMessager(t *testing.T) {
	// --- Test Notef ---
	t.Run("Notef", func(t *testing.T) {
		var buf bytes.Buffer
		m := NewMessager(&buf, false, true, false) // Not silent, trace enabled
		m.Notef("this is a note")
		if !strings.Contains(buf.String(), "Note: this is a note") {
			t.Errorf("Notef should print when trace is enabled, but output was: %q", buf.String())
		}

		buf.Reset()
		m = NewMessager(&buf, false, false, false) // Trace disabled
		m.Notef("this should not print")
		if buf.String() != "" {
			t.Errorf("Notef should not print when trace is disabled, but output was: %q", buf.String())
		}
	})

	// --- Test Warnf ---
	t.Run("Warnf", func(t *testing.T) {
		var buf bytes.Buffer
		m := NewMessager(&buf, false, false, false) // Not silent
		m.Warnf("this is a warning")
		if !strings.Contains(buf.String(), "Warning: this is a warning") {
			t.Errorf("Warnf should print when not silent, but output was: %q", buf.String())
		}

		buf.Reset()
		m = NewMessager(&buf, true, false, false) // Silent
		m.Warnf("this should not print")
		if buf.String() != "" {
			t.Errorf("Warnf should not print when silent, but output was: %q", buf.String())
		}
	})

	// --- Test Errorf ---
	t.Run("Errorf", func(t *testing.T) {
		var buf bytes.Buffer
		m := NewMessager(&buf, true, false, false) // Silent, no show-error
		m.Errorf("this should not print")
		if buf.String() != "" {
			t.Errorf("Errorf should not print when silent, but output was: %q", buf.String())
		}

		buf.Reset()
		m = NewMessager(&buf, true, false, true) // Silent, but with show-error
		m.Errorf("this should print")
		if !strings.Contains(buf.String(), "curl: this should print") {
			t.Errorf("Errorf should print when show-error is true, but output was: %q", buf.String())
		}
	})

	// --- Test Helpf ---
	t.Run("Helpf", func(t *testing.T) {
		var buf bytes.Buffer
		m := NewMessager(&buf, false, false, false)
		m.Helpf("invalid option")
		output := buf.String()
		if !strings.Contains(output, "curl: invalid option") {
			t.Error("Helpf did not print the custom message")
		}
		if !strings.Contains(output, "try 'curl --help'") {
			t.Error("Helpf did not print the standard trailer")
		}
	})

	// --- Test Word Wrapping ---
	t.Run("Word Wrapping", func(t *testing.T) {
		// Set a small terminal width for testing.
		t.Setenv("COLUMNS", "40")

		var buf bytes.Buffer
		m := NewMessager(&buf, false, false, false) // Not silent
		longMessage := "This is a very long line of text that should definitely be wrapped by the function."
		m.Warnf("%s", longMessage)

		output := buf.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")

		if len(lines) <= 1 {
			t.Errorf("Word wrapping failed; expected multiple lines, got 1. Output:\n%s", output)
		}

		// Check that the prefix is on each line.
		for _, line := range lines {
			if !strings.HasPrefix(line, "Warning: ") {
				t.Errorf("Wrapped line does not have prefix: %q", line)
			}
		}
	})
}
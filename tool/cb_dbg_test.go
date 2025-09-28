package tool

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestTracer(t *testing.T) {
	t.Run("TracePlain", func(t *testing.T) {
		var buf bytes.Buffer
		tracer := NewTracer(&buf, TracePlain, false)

		// Test header out
		tracer.Trace(InfoTypeHeaderOut, []byte("GET / HTTP/1.1\r\n"))
		if !strings.Contains(buf.String(), "> GET / HTTP/1.1") {
			t.Errorf("Plain trace for HeaderOut is incorrect. Got:\n%s", buf.String())
		}
		buf.Reset()

		// Test data in
		tracer.Trace(InfoTypeDataIn, []byte{1, 2, 3, 4, 5})
		if !strings.Contains(buf.String(), "[5 bytes data]") {
			t.Errorf("Plain trace for DataIn is incorrect. Got:\n%s", buf.String())
		}
	})

	t.Run("TraceBin", func(t *testing.T) {
		var buf bytes.Buffer
		tracer := NewTracer(&buf, TraceBin, false)
		data := []byte("Hello\x00World")
		tracer.Trace(InfoTypeDataOut, data)
		output := buf.String()

		if !strings.Contains(output, "=> Send data, 11 bytes (0xb)") {
			t.Error("Bin trace header is incorrect")
		}
		if !strings.Contains(output, "48 65 6c 6c 6f") { // "Hello" in hex
			t.Error("Bin trace hex part is incorrect")
		}
		if !strings.Contains(output, "Hello.World") { // Non-printable is replaced by '.'
			t.Error("Bin trace ASCII part is incorrect")
		}
	})

	t.Run("TraceASCII", func(t *testing.T) {
		var buf bytes.Buffer
		tracer := NewTracer(&buf, TraceASCII, false)
		data := []byte("Just ASCII")
		tracer.Trace(InfoTypeDataIn, data)
		output := buf.String()

		if !strings.Contains(output, "<= Recv data, 10 bytes (0xa)") {
			t.Error("ASCII trace header is incorrect")
		}
		// Should NOT contain the hex part
		if strings.Contains(output, "4a 75 73 74") {
			t.Error("ASCII trace should not contain hex dump")
		}
		if !strings.Contains(output, "Just ASCII") {
			t.Error("ASCII trace ASCII part is incorrect")
		}
	})

	t.Run("TraceTime", func(t *testing.T) {
		var buf bytes.Buffer
		tracer := NewTracer(&buf, TracePlain, true)
		tracer.Trace(InfoTypeHeaderOut, []byte("Test\n"))
		output := buf.String()

		// Check for a timestamp like "HH:MM:SS.ffffff "
		re := regexp.MustCompile(`^\d{2}:\d{2}:\d{2}\.\d{6} `)
		if !re.MatchString(output) {
			t.Errorf("Trace output with time enabled is missing timestamp. Got:\n%s", output)
		}
	})
}
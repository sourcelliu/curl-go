package tool

import (
	"fmt"
	"io"
	"time"
)

// InfoType is a translation of the C enum `curl_infotype`.
type InfoType int

const (
	InfoTypeText InfoType = iota
	InfoTypeHeaderIn
	InfoTypeHeaderOut
	InfoTypeDataIn
	InfoTypeDataOut
	InfoTypeSSLDataIn
	InfoTypeSSLDataOut
)

// TraceMode defines the type of trace output.
type TraceMode int

const (
	TraceNone TraceMode = iota
	TraceBin
	TraceASCII
	TracePlain
)

// Tracer handles printing of debug/trace data.
// It encapsulates the state and logic from `tool_cb_dbg.c`.
type Tracer struct {
	Writer    io.Writer
	Mode      TraceMode
	TraceTime bool
}

// NewTracer creates a new tracer.
func NewTracer(writer io.Writer, mode TraceMode, traceTime bool) *Tracer {
	return &Tracer{Writer: writer, Mode: mode, TraceTime: traceTime}
}

// Trace is the Go equivalent of the `tool_debug_cb` C function.
// It receives trace info and data and prints it according to the configured mode.
func (tr *Tracer) Trace(infoType InfoType, data []byte) {
	if tr.Mode == TraceNone || tr.Writer == nil {
		return
	}

	var timebuf string
	if tr.TraceTime {
		now := time.Now()
		timebuf = fmt.Sprintf("%02d:%02d:%02d.%06d ",
			now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000)
	}

	if tr.Mode == TracePlain {
		tr.tracePlain(infoType, data, timebuf)
		return
	}

	var text string
	switch infoType {
	case InfoTypeText:
		fmt.Fprintf(tr.Writer, "%s== Info: %s", timebuf, string(data))
		return
	case InfoTypeHeaderIn:
		text = "<= Recv header"
	case InfoTypeHeaderOut:
		text = "=> Send header"
	case InfoTypeDataIn:
		text = "<= Recv data"
	case InfoTypeDataOut:
		text = "=> Send data"
	case InfoTypeSSLDataIn:
		text = "<= Recv SSL data"
	case InfoTypeSSLDataOut:
		text = "=> Send SSL data"
	default:
		return
	}

	dump(tr.Writer, timebuf, text, data, tr.Mode)
}

// tracePlain handles the logic for the simple "-v" style trace.
func (tr *Tracer) tracePlain(infoType InfoType, data []byte, timebuf string) {
	prefix := "* "
	switch infoType {
	case InfoTypeHeaderIn:
		prefix = "< "
	case InfoTypeHeaderOut:
		prefix = "> "
	}

	// For TEXT, HEADER_IN, HEADER_OUT, just print the line with a prefix.
	// We simplify the C code's complex state management for newlines.
	if infoType == InfoTypeText || infoType == InfoTypeHeaderIn || infoType == InfoTypeHeaderOut {
		fmt.Fprintf(tr.Writer, "%s%s%s", timebuf, prefix, string(data))
		// Ensure it ends with a newline
		if len(data) == 0 || data[len(data)-1] != '\n' {
			fmt.Fprintln(tr.Writer)
		}
	}
	// For data types, the C code prints a summary. We do the same.
	if infoType == InfoTypeDataIn || infoType == InfoTypeDataOut ||
		infoType == InfoTypeSSLDataIn || infoType == InfoTypeSSLDataOut {
		fmt.Fprintf(tr.Writer, "%s%s[%d bytes data]\n", timebuf, prefix, len(data))
	}
}

// dump creates a hex dump of the data. It's a translation of the C `dump` function.
func dump(writer io.Writer, timebuf, text string, data []byte, mode TraceMode) {
	fmt.Fprintf(writer, "%s%s, %d bytes (0x%x)\n", timebuf, text, len(data), len(data))

	width := 16 // 0x10, for TRACE_BIN
	if mode == TraceASCII {
		width = 64 // 0x40
	}

	for i := 0; i < len(data); i += width {
		fmt.Fprintf(writer, "%04x: ", i)

		// Hex dump
		if mode == TraceBin {
			for c := 0; c < width; c++ {
				if i+c < len(data) {
					fmt.Fprintf(writer, "%02x ", data[i+c])
				} else {
					fmt.Fprint(writer, "   ")
				}
			}
		}

		// ASCII dump
		for c := 0; c < width && i+c < len(data); c++ {
			char := data[i+c]
			if char >= 0x20 && char < 0x7f {
				writer.Write([]byte{char})
			} else {
				writer.Write([]byte{'.'})
			}
		}
		fmt.Fprintln(writer)
	}
}
package tool

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestFormatBytes(t *testing.T) {
	testCases := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"less than 100k", 12345, "12345"},
		{"kilobytes", 1234567, "1205k"},
		{"megabytes with decimal", 12345678, "11.7M"},
		{"megabytes without decimal", 1234567890, "1177M"},
		{"gigabytes with decimal", 12345678901, "11.4G"}, // C uses integer math, which truncates 11.498 to 11.4
		{"gigabytes without decimal", 123456789012, " 114G"},
		{"terabytes", 12345678901234, "  11T"},
		{"petabytes", 12345678901234567, "  10P"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formatBytes(tc.bytes)
			// Pad the expected result to 5 characters for consistent comparison
			paddedExpected := "     " + tc.expected
			paddedExpected = paddedExpected[len(paddedExpected)-5:]
			if result != paddedExpected {
				t.Errorf("formatBytes(%d) = %q; want %q", tc.bytes, result, paddedExpected)
			}
		})
	}
}

func TestFormatSeconds(t *testing.T) {
	testCases := []struct {
		name     string
		seconds  int64
		expected string
	}{
		{"zero", 0, "--:--:--"},
		{"simple time", 3723, "01:02:03"}, // 1h 2m 3s
		{"less than 100 hours", 359999, "99:59:59"},
		{"more than 100 hours", 360000, "  4d 04h"},
		{"large days", 8640000, "100d 00h"}, // C code shows hours even if zero
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if result := formatSeconds(tc.seconds); result != tc.expected {
				t.Errorf("formatSeconds(%d) = %q; want %q", tc.seconds, result, tc.expected)
			}
		})
	}
}

func TestProgressBar(t *testing.T) {
	var buf bytes.Buffer
	// Make the start time deterministic for consistent output
	startTime := time.Now().Add(-10 * time.Second)
	p := &ProgressBar{
		Writer:    &buf,
		StartTime: startTime,
	}

	stats := []TransferStats{
		{DLTotal: 2048, DLNow: 1024, ULTotal: 0, ULNow: 0},
	}

	p.Render(stats, false)
	output := buf.String()

	// Check for key components of the progress bar
	if !strings.Contains(output, "DL% UL%") {
		t.Error("Output does not contain header")
	}
	if !strings.Contains(output, "50") { // 1024 is 50% of 2048
		t.Error("Output does not contain correct DL percentage")
	}
	if !strings.Contains(output, "1024") { // Downloaded bytes
		t.Error("Output does not contain correct downloaded bytes")
	}
	if !strings.Contains(output, "00:00:10") { // 10 seconds spent
		t.Error("Output does not contain correct time spent")
	}

	// Test final render adds a newline
	buf.Reset()
	p.Render(stats, true)
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Error("Final render should end with a newline")
	}
}
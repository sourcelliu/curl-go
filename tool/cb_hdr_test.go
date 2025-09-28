package tool

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestParseContentDispositionFilename(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"simple filename", `attachment; filename="foo.html"`, "foo.html", false},
		{"unquoted filename", `attachment; filename=bar.jpg`, "bar.jpg", false},
		{"form-data with filename", `form-data; name="field"; filename="file.zip"`, "file.zip", false},
		{"no filename", `attachment`, "", false},
		{"empty filename", `attachment; filename=""`, "", false},
		{"invalid format", `attachment; filename=`, "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename, err := parseContentDispositionFilename(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("parseContentDispositionFilename() error = %v, wantErr %v", err, tc.wantErr)
			}
			if filename != tc.expected {
				t.Errorf("filename = %q; want %q", filename, tc.expected)
			}
		})
	}
}

func TestHeaderProcessor(t *testing.T) {
	t.Run("header dumping", func(t *testing.T) {
		var buf1, buf2 bytes.Buffer
		hp := NewHeaderProcessor()
		hp.HeaderWriters = []io.Writer{&buf1, &buf2}
		headerLine := "Content-Type: application/json\r\n"

		if err := hp.Process(headerLine); err != nil {
			t.Fatalf("Process() failed: %v", err)
		}

		if buf1.String() != headerLine {
			t.Errorf("Buffer 1 content = %q; want %q", buf1.String(), headerLine)
		}
		if buf2.String() != headerLine {
			t.Errorf("Buffer 2 content = %q; want %q", buf2.String(), headerLine)
		}
	})

	t.Run("etag saving", func(t *testing.T) {
		var etagBuf bytes.Buffer
		hp := NewHeaderProcessor()
		hp.ETagWriter = &etagBuf
		headerLine := `ETag: "12345-abcde"` + "\r\n"

		if err := hp.Process(headerLine); err != nil {
			t.Fatalf("Process() failed: %v", err)
		}

		expected := `"12345-abcde"` + "\n"
		if etagBuf.String() != expected {
			t.Errorf("ETag buffer = %q; want %q", etagBuf.String(), expected)
		}
	})

	t.Run("content disposition parsing", func(t *testing.T) {
		hp := NewHeaderProcessor()
		hp.HonorContentDisposition = true
		headerLine := `Content-Disposition: attachment; filename="download.zip"` + "\r\n"

		if err := hp.Process(headerLine); err != nil {
			t.Fatalf("Process() failed: %v", err)
		}

		if hp.FilenameFromDisposition != "download.zip" {
			t.Errorf("FilenameFromDisposition = %q; want %q", hp.FilenameFromDisposition, "download.zip")
		}
		if hp.HonorContentDisposition {
			t.Error("HonorContentDisposition should be false after finding a filename")
		}
	})

	t.Run("combined functionality", func(t *testing.T) {
		var dumpBuf, etagBuf bytes.Buffer
		hp := NewHeaderProcessor()
		hp.HeaderWriters = []io.Writer{&dumpBuf}
		hp.ETagWriter = &etagBuf
		hp.HonorContentDisposition = true

		headers := []string{
			"HTTP/1.1 200 OK\r\n",
			`Content-Disposition: inline; filename="display.jpg"` + "\r\n",
			`ETag: "etag-value"` + "\r\n",
			"\r\n",
		}

		for _, h := range headers {
			if err := hp.Process(h); err != nil {
				t.Fatalf("Process() failed on header %q: %v", h, err)
			}
		}

		// Verify all parts
		if !strings.Contains(dumpBuf.String(), "HTTP/1.1 200 OK") {
			t.Error("Dump buffer is missing headers")
		}
		if etagBuf.String() != `"etag-value"`+"\n" {
			t.Errorf("ETag not saved correctly, got %q", etagBuf.String())
		}
		if hp.FilenameFromDisposition != "display.jpg" {
			t.Errorf("Filename not extracted correctly, got %q", hp.FilenameFromDisposition)
		}
	})
}
package tool

import (
	"reflect"
	"testing"
)

func TestParseFormString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []*FormPart
		wantErr  bool
	}{
		{
			name:  "simple literal",
			input: "name=value",
			expected: []*FormPart{
				{Name: "name", Value: "value", Type: PartTypeLiteral, Filename: "value"},
			},
		},
		{
			name:  "file upload",
			input: "file=@path/to/file.txt",
			expected: []*FormPart{
				{Name: "file", Value: "path/to/file.txt", Type: PartTypeFile, Filename: "path/to/file.txt"},
			},
		},
		{
			name:  "file content",
			input: "data=<path/to/data.txt",
			expected: []*FormPart{
				{Name: "data", Value: "path/to/data.txt", Type: PartTypeDataFile, Filename: "path/to/data.txt"},
			},
		},
		{
			name:  "file upload with type and filename",
			input: `upload=@file.zip;type=application/zip;filename=archive.zip`,
			expected: []*FormPart{
				{Name: "upload", Value: "file.zip", Type: PartTypeFile, ContentType: "application/zip", Filename: "archive.zip"},
			},
		},
		{
			name:  "quoted filename with spaces",
			input: `attachment=@"my document.pdf";type="application/pdf"`,
			expected: []*FormPart{
				{Name: "attachment", Value: "my document.pdf", Type: PartTypeFile, ContentType: "application/pdf", Filename: "my document.pdf"},
			},
		},
		{
			name:  "multiple files",
			input: "images=@img1.jpg,@img2.png",
			expected: []*FormPart{
				{Name: "images", Value: "img1.jpg", Type: PartTypeFile, Filename: "img1.jpg"},
				{Name: "images", Value: "img2.png", Type: PartTypeFile, Filename: "img2.png"},
			},
		},
		{
			name:  "complex multiple files",
			input: `assets=@"a.txt";type=text/plain,@"b.zip";filename="b archive.zip"`,
			expected: []*FormPart{
				{Name: "assets", Value: "a.txt", Type: PartTypeFile, ContentType: "text/plain", Filename: "a.txt"},
				{Name: "assets", Value: "b.zip", Type: PartTypeFile, Filename: "b archive.zip"},
			},
		},
		{
			name:    "missing equals",
			input:   "namevalue",
			wantErr: true,
		},
		{
			name:    "missing name",
			input:   "=value",
			wantErr: true,
		},
		{
			name:    "literal with comma",
			input:   "name=value,another",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseFormString(tc.input)

			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseFormString() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				// For better diffing in test logs
				for i := 0; i < len(result) || i < len(tc.expected); i++ {
					var r, e interface{}
					if i < len(result) {
						r = result[i]
					}
					if i < len(tc.expected) {
						e = tc.expected[i]
					}
					if !reflect.DeepEqual(r, e) {
						t.Errorf("Mismatch at index %d:\nGot:    %+v\nWanted: %+v", i, r, e)
					}
				}
				t.Fatalf("Overall result mismatch")
			}
		})
	}
}
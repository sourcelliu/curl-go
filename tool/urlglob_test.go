package tool

import (
	"reflect"
	"sort"
	"testing"
)

func TestExpandURLGlob(t *testing.T) {
	testCases := []struct {
		name     string
		pattern  string
		expected []string
		wantErr  bool
	}{
		{
			name:     "no globbing",
			pattern:  "http://example.com/page.html",
			expected: []string{"http://example.com/page.html"},
		},
		{
			name:     "simple set glob",
			pattern:  "http://example.{one,two,three}.com",
			expected: []string{"http://example.one.com", "http://example.two.com", "http://example.three.com"},
		},
		{
			name:     "numeric range glob",
			pattern:  "http://example.com/page[1-4].html",
			expected: []string{"http://example.com/page1.html", "http://example.com/page2.html", "http://example.com/page3.html", "http://example.com/page4.html"},
		},
		{
			name:     "character range glob",
			pattern:  "ftp://ftp.example.com/file[a-c].txt",
			expected: []string{"ftp://ftp.example.com/filea.txt", "ftp://ftp.example.com/fileb.txt", "ftp://ftp.example.com/filec.txt"},
		},
		{
			name:     "multiple globs",
			pattern:  "http://{site,host}.example.com/page[1-2]",
			expected: []string{"http://site.example.com/page1", "http://site.example.com/page2", "http://host.example.com/page1", "http://host.example.com/page2"},
		},
		{
			name:     "escaped brace",
			pattern:  `http://example.com/\{a,b}`,
			expected: []string{`http://example.com/{a,b}`},
		},
		{
			name:     "escaped bracket",
			pattern:  `http://example.com/\[1-3]`,
			expected: []string{`http://example.com/[1-3]`},
		},
		{
			name:    "unmatched brace",
			pattern: "http://example.com/{a,b",
			wantErr: true,
		},
		{
			name:    "unmatched bracket",
			pattern: "http://example.com/[1-5",
			wantErr: true,
		},
		{
			name:    "invalid numeric range",
			pattern: "http://example.com/[5-1]",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExpandURLGlob(tc.pattern)

			if (err != nil) != tc.wantErr {
				t.Fatalf("ExpandURLGlob() error = %v, wantErr %v", err, tc.wantErr)
			}

			// Sort both slices for stable comparison
			sort.Strings(result)
			sort.Strings(tc.expected)

			if !tc.wantErr && !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("ExpandURLGlob() = %v; want %v", result, tc.expected)
			}
		})
	}
}
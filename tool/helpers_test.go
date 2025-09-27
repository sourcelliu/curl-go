package tool

import (
	"strings"
	"testing"
)

func TestParameterError_String(t *testing.T) {
	testCases := []struct {
		name     string
		err      ParameterError
		expected string
	}{
		{
			name:     "ok",
			err:      ParamOK,
			expected: "ok",
		},
		{
			name:     "unknown option",
			err:      ParamOptionUnknown,
			expected: "is unknown",
		},
		{
			name:     "out of memory",
			err:      ParamNoMem,
			expected: "out of memory",
		},
		{
			name:     "unknown error code",
			err:      ParameterError(999),
			expected: "unknown error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if result := tc.err.String(); result != tc.expected {
				t.Errorf("String() = %q; want %q", result, tc.expected)
			}
		})
	}
}

func TestHTTPRequestManager(t *testing.T) {
	t.Run("set once", func(t *testing.T) {
		var m HTTPRequestManager
		err := m.Set(HTTPRequestGet)
		if err != nil {
			t.Errorf("Set() returned an unexpected error: %v", err)
		}
		if m.request != HTTPRequestGet {
			t.Errorf("request = %v; want %v", m.request, HTTPRequestGet)
		}
	})

	t.Run("set same twice", func(t *testing.T) {
		var m HTTPRequestManager
		m.Set(HTTPRequestHead)
		err := m.Set(HTTPRequestHead)
		if err != nil {
			t.Errorf("Set() returned an unexpected error on second call: %v", err)
		}
	})

	t.Run("set conflicting", func(t *testing.T) {
		var m HTTPRequestManager
		m.Set(HTTPRequestGet)
		err := m.Set(HTTPRequestPut)
		if err == nil {
			t.Error("Set() did not return an error on conflicting request")
		} else {
			// Check if the error message is as expected.
			expectedMsg := "you can only select one HTTP request method"
			if !strings.Contains(err.Error(), expectedMsg) {
				t.Errorf("Set() error message = %q; want to contain %q", err.Error(), expectedMsg)
			}
		}
	})
}

func TestCustomRequestHelper(t *testing.T) {
	testCases := []struct {
		name     string
		req      HTTPRequest
		method   string
		expected string
	}{
		{
			name:     "no custom method",
			req:      HTTPRequestGet,
			method:   "",
			expected: "",
		},
		{
			name:     "unnecessary GET",
			req:      HTTPRequestGet,
			method:   "GET",
			expected: "Unnecessary use of -X or --request",
		},
		{
			name:     "unnecessary HEAD",
			req:      HTTPRequestHead,
			method:   "HEAD",
			expected: "Unnecessary use of -X or --request",
		},
		{
			name:     "custom HEAD warning",
			req:      HTTPRequestGet,
			method:   "HEAD",
			expected: "Setting custom HTTP method to HEAD with -X/--request",
		},
		{
			name:     "valid custom method",
			req:      HTTPRequestGet,
			method:   "DELETE",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CustomRequestHelper(tc.req, tc.method)
			if !strings.Contains(result, tc.expected) {
				t.Errorf("CustomRequestHelper() = %q; want to contain %q", result, tc.expected)
			}
		})
	}
}
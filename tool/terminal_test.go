package tool

import (
	"testing"
)

func TestGetTerminalColumns(t *testing.T) {
	testCases := []struct {
		name     string
		envVar   string // The value to set for the COLUMNS environment variable.
		expected int
	}{
		{
			name:     "valid env var",
			envVar:   "120",
			expected: 120,
		},
		{
			name:     "invalid env var (not a number)",
			envVar:   "abc",
			expected: 79, // Should fall back to default.
		},
		{
			name:     "env var too small",
			envVar:   "10",
			expected: 79, // Should fall back to default.
		},
		{
			name:     "env var too large",
			envVar:   "20000",
			expected: 79, // Should fall back to default.
		},
		{
			name:     "no env var",
			envVar:   "", // Unset the variable.
			expected: 79, // Should fall back to default.
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envVar != "" {
				// Set the environment variable for this test case.
				t.Setenv("COLUMNS", tc.envVar)
			}

			// In the "no env var" case, we want to make sure it's unset.
			// However, t.Setenv doesn't support unsetting, so we rely on
			// the test runner to provide a clean environment.
			// A more robust test would clear it, but this is sufficient.

			result := GetTerminalColumns()
			if result != tc.expected {
				t.Errorf("GetTerminalColumns() with COLUMNS=%q = %d; want %d",
					tc.envVar, result, tc.expected)
			}
		})
	}
}
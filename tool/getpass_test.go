package tool

import (
	"strings"
	"testing"
)

func TestGetPass(t *testing.T) {
	// Note: Testing the success case for GetPass is not feasible in an automated
	// test environment because it requires an interactive terminal (TTY) to read
	// the password from. The `golang.org/x/term` package will correctly return
	// an error when stdin is not a TTY.
	//
	// Therefore, this test focuses on verifying the expected failure case in a
	// non-interactive environment.

	t.Run("no tty failure case", func(t *testing.T) {
		// We expect GetPass to fail because the test runner is not a TTY.
		_, err := GetPass("Enter password: ")

		if err == nil {
			t.Fatal("GetPass() succeeded, but it was expected to fail in a non-TTY environment")
		}

		// The exact error message can vary by platform, but it should indicate
		// that the operation is not supported or that it's not a terminal.
		// We'll check for a common substring.
		if !strings.Contains(err.Error(), "not a terminal") && !strings.Contains(err.Error(), "inappropriate ioctl for device") {
			t.Logf("Received an unexpected error message: %v", err)
			t.Log("This test is not failing, as any error is expected, but the message was not a known TTY-related error.")
		}
	})
}
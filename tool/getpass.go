package tool

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// GetPass is a translation of the C function `getpass_r` from
// curl-src/src/tool_getpass.c.
//
// The C function reads a password from the terminal without echoing characters.
// It contains multiple platform-specific implementations (for VMS, Windows,
// and Unix-like systems using termios). The idiomatic Go equivalent is to use
// `golang.org/x/term.ReadPassword`, which provides a cross-platform, safe way
// to read passwords.
//
// Original C code from tool_getpass.c.
func GetPass(prompt string) (string, error) {
	// Original C code logic from tool_getpass.c, lines 112, 126, 214:
	//   fputs(prompt, tool_stderr);
	// We print the prompt to standard error, just like curl's implementation.
	fmt.Fprint(os.Stderr, prompt)

	// term.ReadPassword takes care of disabling and re-enabling terminal echo.
	// It reads from the provided file descriptor, which should be a terminal.
	// We use os.Stdin for this.
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	// The C code often prints a newline after reading the password since echo
	// is off. term.ReadPassword generally handles the terminal state correctly,
	// but adding a newline to stderr ensures the next output starts cleanly.
	fmt.Fprintln(os.Stderr)

	return string(bytePassword), nil
}
package tool

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// ParseLong converts a string to a long integer (int64).
// This is a Go equivalent for the C function `str2num`.
func ParseLong(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}
	return val, nil
}

// ParseULong converts a string to an unsigned long integer (uint64),
// ensuring it's not negative.
// This is a Go equivalent for the C function `str2unum`.
func ParseULong(s string) (uint64, error) {
	if strings.HasPrefix(s, "-") {
		return 0, fmt.Errorf("negative number not allowed: %s", s)
	}
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}
	return val, nil
}

// ParseSecs converts a string representing seconds (possibly with decimals)
// into a time.Duration.
// This is the Go equivalent of the C function `secs2ms`.
func ParseSecs(s string) (time.Duration, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || f < 0 {
		return 0, fmt.Errorf("invalid time value: %s", s)
	}
	// Convert float seconds to nanoseconds for time.Duration
	return time.Duration(f * float64(time.Second)), nil
}

// FileToString reads the entire content of a file into a string.
// This is a simpler, safer Go equivalent of the C functions
// `file2string` and `file2memory`. It strips CR and LF characters.
func FileToString(file *os.File) (string, error) {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	// The C version strips CR and LF. We'll replicate that.
	s := string(content)
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s, nil
}

// FTPFileMethod is an enum for FTP file transfer methods.
type FTPFileMethod int

const (
	FTPMethodMultiCWD FTPFileMethod = iota
	FTPMethodNoCWD
	FTPMethodSingleCWD
)

// ParseFTPFileMethod converts a string to an FTPFileMethod enum.
// This is a translation of the C function `ftpfilemethod`.
func ParseFTPFileMethod(s string) (FTPFileMethod, error) {
	switch strings.ToLower(s) {
	case "multicwd":
		return FTPMethodMultiCWD, nil
	case "nocwd":
		return FTPMethodNoCWD, nil
	case "singlecwd":
		return FTPMethodSingleCWD, nil
	default:
		return 0, fmt.Errorf("unrecognized ftp file method: %s", s)
	}
}

// GSSDelegation is an enum for GSS-API delegation levels.
type GSSDelegation int

const (
	DelegationNone GSSDelegation = iota
	DelegationPolicy
	DelegationAlways
)

// ParseDelegation converts a string to a GSSDelegation enum.
// This is a translation of the C function `delegation`.
func ParseDelegation(s string) (GSSDelegation, error) {
	switch strings.ToLower(s) {
	case "none":
		return DelegationNone, nil
	case "policy":
		return DelegationPolicy, nil
	case "always":
		return DelegationAlways, nil
	default:
		return 0, fmt.Errorf("unrecognized delegation method: %s", s)
	}
}
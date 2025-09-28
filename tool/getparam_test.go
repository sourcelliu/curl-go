package tool

import (
	"runtime"
	"testing"
	"time"
)

func TestParameterParser_ParseOne(t *testing.T) {
	testCases := []struct {
		name        string
		flag        string
		nextArg     string
		expectUsed  bool
		expectErr   bool
		setupConfig func(config *OperationConfig) // Optional setup
		checkConfig func(t *testing.T, config *OperationConfig)
	}{
		{
			name:       "long option with argument",
			flag:       "--user-agent",
			nextArg:    "test-agent/1.0",
			expectUsed: true,
			checkConfig: func(t *testing.T, config *OperationConfig) {
				if config.UserAgent != "test-agent/1.0" {
					t.Errorf("UserAgent = %q; want %q", config.UserAgent, "test-agent/1.0")
				}
			},
		},
		{
			name:       "short option with bundled argument",
			flag:       "-odata.txt",
			expectUsed: false,
			checkConfig: func(t *testing.T, config *OperationConfig) {
				if len(config.URLList) != 1 || config.URLList[0].Outfile != "data.txt" {
					t.Errorf("Outfile was not set correctly for bundled arg")
				}
			},
		},
		{
			name:      "unknown option",
			flag:      "--non-existent-flag",
			expectErr: true,
		},
		{
			name:      "option requires argument but none given",
			flag:      "--output",
			nextArg:   "",
			expectErr: true,
		},
		{
			name:       "valid range",
			flag:       "--range",
			nextArg:    "0-1023",
			expectUsed: true,
			checkConfig: func(t *testing.T, config *OperationConfig) {
				if config.Range != "0-1023" {
					t.Errorf("Range = %q; want %q", config.Range, "0-1023")
				}
			},
		},
		{
			name:      "range with no dash",
			flag:      "--range",
			nextArg:   "1024",
			expectErr: true,
		},
		{
			name:    "range conflicts with continue-at",
			flag:    "--range",
			nextArg: "0-",
			setupConfig: func(config *OperationConfig) {
				config.UseResume = true
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			global := NewGlobalConfig()
			parser := NewParameterParser(global)
			if tc.setupConfig != nil {
				tc.setupConfig(global.Last)
			}

			used, err := parser.ParseOne(tc.flag, tc.nextArg)

			if (err != nil) != tc.expectErr {
				t.Fatalf("ParseOne() error = %v, wantErr %v", err, tc.expectErr)
			}
			if used != tc.expectUsed {
				t.Errorf("ParseOne() used = %v; want %v", used, tc.expectUsed)
			}

			if tc.checkConfig != nil {
				tc.checkConfig(t, global.Last)
			}
		})
	}
}
func TestParameterParser_Parse(t *testing.T) {
	t.Run("full command line", func(t *testing.T) {
		args := []string{
			"-v",
			"--user-agent", "test-agent",
			"http://example.com",
			"-H", "X-Test: true",
			"--connect-timeout", "2.5",
			"-f",
		}
		global := NewGlobalConfig()
		parser := NewParameterParser(global)

		err := parser.Parse(args)
		if err != nil {
			t.Fatalf("Parse() failed: %v", err)
		}

		config := global.Last
		if config.UserAgent != "test-agent" {
			t.Errorf("UserAgent = %q; want %q", config.UserAgent, "test-agent")
		}
		if len(config.URLList) != 1 || config.URLList[0].URL != "http://example.com" {
			t.Errorf("URL not parsed correctly")
		}
		if len(config.Headers) != 1 || config.Headers[0] != "X-Test: true" {
			t.Errorf("Header not parsed correctly")
		}
		if config.ConnectTimeout != 2500*time.Millisecond {
			t.Errorf("ConnectTimeout = %v; want %v", config.ConnectTimeout, 2500*time.Millisecond)
		}
		if !config.FailOnError {
			t.Error("FailOnError should be true")
		}
	})

	t.Run("end of flags", func(t *testing.T) {
		args := []string{"--", "-this-is-a-url"}
		global := NewGlobalConfig()
		parser := NewParameterParser(global)
		if err := parser.Parse(args); err != nil {
			t.Fatalf("Parse() failed: %v", err)
		}
		config := global.Last
		if len(config.URLList) != 1 || config.URLList[0].URL != "-this-is-a-url" {
			t.Errorf("URL after -- was not parsed correctly")
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		args := []string{"--verbose", "--non-existent"}
		global := NewGlobalConfig()
		parser := NewParameterParser(global)
		err := parser.Parse(args)
		if err == nil {
			t.Error("Parse() did not return an error for an invalid flag")
		}
	})
}

func TestParameterParser_Auth(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected AuthType
	}{
		{
			name:     "basic only",
			args:     []string{"--basic"},
			expected: AuthBasic,
		},
		{
			name:     "digest only",
			args:     []string{"--digest"},
			expected: AuthDigest,
		},
		{
			name:     "multiple auth flags",
			args:     []string{"--basic", "--ntlm"},
			expected: AuthBasic | AuthNTLM,
		},
		{
			name:     "anyauth overrides",
			args:     []string{"--basic", "--anyauth", "--digest"},
			expected: AuthAny,
		},
		{
			name:     "anyauth alone",
			args:     []string{"--anyauth"},
			expected: AuthAny,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			global := NewGlobalConfig()
			parser := NewParameterParser(global)
			if err := parser.Parse(tc.args); err != nil {
				t.Fatalf("Parse() failed: %v", err)
			}
			if global.Last.AuthType != uint(tc.expected) {
				t.Errorf("AuthType = %d; want %d", global.Last.AuthType, tc.expected)
			}
		})
	}
}

func TestNewGlobalConfig(t *testing.T) {
	g := NewGlobalConfig()

	if g == nil {
		t.Fatal("NewGlobalConfig() returned nil")
	}
	if g.First == nil {
		t.Error("GlobalConfig.First should not be nil")
	}
	if g.Last == nil {
		t.Error("GlobalConfig.Last should not be nil")
	}
	if g.First != g.Last {
		t.Error("GlobalConfig.First and .Last should point to the same initial config")
	}
}

func TestParseCertParameter(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedCert string
		expectedPass string
		onlyOnOS     string // "windows" or ""
	}{
		{"simple case", "mycert:mypass", "mycert", "mypass", ""},
		{"no password", "mycert", "mycert", "", ""},
		{"pkcs11 uri", "pkcs11:token=my-token;object=my-cert", "pkcs11:token=my-token;object=my-cert", "", ""},
		{"escaped colon", `my\:cert:mypass`, `my:cert`, "mypass", ""},
		{"windows path", `C:\path\to\cert.pem`, `C:\path\to\cert.pem`, "", "windows"},
		{"windows path with password", `C:\path:password`, `C:\path`, "password", "windows"},
		{"empty string", "", "", "", ""},
		{"password only", ":mypass", "", "mypass", ""},
		{"cert only", "mycert:", "mycert", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.onlyOnOS != "" && tc.onlyOnOS != runtime.GOOS {
				t.Skipf("Skipping test case %q on OS %q", tc.name, runtime.GOOS)
			}

			cert, pass := parseCertParameter(tc.input)
			if cert != tc.expectedCert {
				t.Errorf("Cert name = %q; want %q", cert, tc.expectedCert)
			}
			if pass != tc.expectedPass {
				t.Errorf("Passphrase = %q; want %q", pass, tc.expectedPass)
			}
		})
	}
}

func TestParseSizeParameter(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"bytes plain", "1024", 1024, false},
		{"bytes with B", "2048B", 2048, false},
		{"kilobytes uppercase", "100K", 102400, false},
		{"megabytes lowercase", "2m", 2 * 1024 * 1024, false},
		{"gigabytes mixed case", "1g", 1 * 1024 * 1024 * 1024, false},
		{"empty string", "", 0, true},
		{"invalid number", "abcK", 0, true},
		{"invalid suffix", "100X", 0, true},
		{"overflow", "9999999999999999999G", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseSizeParameter(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("parseSizeParameter(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}
			if !tc.wantErr && result != tc.expected {
				t.Errorf("parseSizeParameter(%q) = %d; want %d", tc.input, result, tc.expected)
			}
		})
	}
}

func TestOptionsMap(t *testing.T) {
	if options == nil {
		t.Fatal("Options map is nil")
	}

	testCases := []struct {
		name      string
		shortName rune
		argType   ArgType
	}{
		{"url", 0, ArgString},
		{"verbose", 'v', ArgBool},
		{"header", 'H', ArgString},
		{"output", 'o', ArgFile},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opt, ok := options[tc.name]
			if !ok {
				t.Fatalf("Option %q not found in map", tc.name)
			}
			if opt.ShortName != tc.shortName {
				t.Errorf("ShortName for %q is %c; want %c", tc.name, opt.ShortName, tc.shortName)
			}
			if opt.Type != tc.argType {
				t.Errorf("ArgType for %q is %v; want %v", tc.name, opt.Type, tc.argType)
			}
		})
	}
}
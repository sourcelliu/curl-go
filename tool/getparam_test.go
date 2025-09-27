package tool

import (
	"testing"
)

func TestParameterParser_ParseOne(t *testing.T) {
	testCases := []struct {
		name        string
		flag        string
		nextArg     string
		expectUsed  bool
		expectErr   bool
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
			name:       "long boolean option",
			flag:       "--insecure",
			expectUsed: false,
			checkConfig: func(t *testing.T, config *OperationConfig) {
				if !config.InsecureOK {
					t.Error("InsecureOK should be true")
				}
			},
		},
		{
			name:       "short option",
			flag:       "-L",
			expectUsed: false,
			checkConfig: func(t *testing.T, config *OperationConfig) {
				if !config.FollowLocation {
					t.Error("FollowLocation should be true")
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
			name:       "unknown option",
			flag:       "--non-existent-flag",
			expectErr:  true,
		},
		{
			name:      "option requires argument but none given",
			flag:      "--output",
			nextArg:   "",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			global := NewGlobalConfig()
			parser := NewParameterParser(global)

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
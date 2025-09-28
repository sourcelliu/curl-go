package tool

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

// ArgType defines the type of argument an option expects.
type ArgType int

const (
	ArgNone ArgType = iota // Stand-alone option
	ArgBool                // Boolean option (e.g., --verbose, --no-verbose)
	ArgString              // Option requires a string argument
	ArgFile                // Option requires a file path argument
)

// Option defines a single command-line option.
type Option struct {
	Name      string
	ShortName rune
	Type      ArgType
	// Handler is the function that applies the option to an OperationConfig.
	Handler func(p *ParameterParser, config *OperationConfig, arg string) error
}

// options is a map of all supported command-line options.
var options = map[string]Option{
	"url":             {Name: "url", Type: ArgString, Handler: handleURL},
	"verbose":         {Name: "verbose", ShortName: 'v', Type: ArgBool, Handler: handleVerbose},
	"header":          {Name: "header", ShortName: 'H', Type: ArgString, Handler: handleHeader},
	"data":            {Name: "data", ShortName: 'd', Type: ArgString, Handler: handleData},
	"request":         {Name: "request", ShortName: 'X', Type: ArgString, Handler: handleString("CustomRequest")},
	"user-agent":      {Name: "user-agent", ShortName: 'A', Type: ArgString, Handler: handleString("UserAgent")},
	"insecure":        {Name: "insecure", ShortName: 'k', Type: ArgBool, Handler: handleBool("InsecureOK")},
	"location":        {Name: "location", ShortName: 'L', Type: ArgBool, Handler: handleBool("FollowLocation")},
	"output":          {Name: "output", ShortName: 'o', Type: ArgFile, Handler: handleOutputFile},
	"remote-name":     {Name: "remote-name", ShortName: 'O', Type: ArgBool, Handler: handleRemoteName},
	"user":            {Name: "user", ShortName: 'u', Type: ArgString, Handler: handleString("UserPassword")},
	"head":            {Name: "head", ShortName: 'I', Type: ArgBool, Handler: handleHead},
	"get":             {Name: "get", ShortName: 'G', Type: ArgBool, Handler: handleBool("UseHTTPGet")},
	"connect-timeout": {Name: "connect-timeout", Type: ArgString, Handler: handleConnectTimeout},
	"fail":            {Name: "fail", ShortName: 'f', Type: ArgBool, Handler: handleBool("FailOnError")},
	"range":           {Name: "range", ShortName: 'r', Type: ArgString, Handler: handleRange},
}

// shortOptions is a reverse map for finding long options by their short name.
var shortOptions = make(map[rune]Option)

func init() {
	for name, opt := range options {
		if opt.ShortName != 0 {
			// Add a reference back to the long name for consistency
			o := opt
			o.Name = name
			shortOptions[opt.ShortName] = o
		}
	}
}

// ParameterParser holds the state for parsing arguments.
type ParameterParser struct {
	Global *GlobalConfig
}

// NewParameterParser creates a new parser.
func NewParameterParser(global *GlobalConfig) *ParameterParser {
	return &ParameterParser{Global: global}
}

// Parse iterates through the command-line arguments and processes them.
func (p *ParameterParser) Parse(args []string) error {
	stillFlags := true
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if stillFlags && strings.HasPrefix(arg, "-") {
			if arg == "--" {
				stillFlags = false // End of flags
				continue
			}

			var nextArg string
			if i+1 < len(args) {
				nextArg = args[i+1]
			}

			usedArg, err := p.ParseOne(arg, nextArg)
			if err != nil {
				return fmt.Errorf("option %s: %w", arg, err)
			}
			if usedArg {
				i++ // The next argument was consumed
			}
		} else {
			// Not a flag, treat as a URL
			_, err := p.ParseOne("--url", arg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ParseOne parses a single flag and its potential argument.
func (p *ParameterParser) ParseOne(flag, nextarg string) (usedArg bool, err error) {
	var opt Option
	var ok bool
	var arg string

	isLongOpt := strings.HasPrefix(flag, "--")
	if isLongOpt {
		optName := strings.TrimPrefix(flag, "--")
		opt, ok = options[optName]
		if !ok {
			return false, fmt.Errorf("unknown option")
		}
	} else { // Short option
		shortName := rune(flag[1])
		opt, ok = shortOptions[shortName]
		if !ok {
			return false, fmt.Errorf("unknown option")
		}
		if len(flag) > 2 {
			arg = flag[2:]
		}
	}

	if opt.Type == ArgString || opt.Type == ArgFile {
		if arg != "" {
			// Argument was bundled
		} else if nextarg != "" {
			arg = nextarg
			usedArg = true
		} else {
			return false, fmt.Errorf("requires an argument")
		}
	}

	if opt.Handler != nil {
		err = opt.Handler(p, p.Global.Last, arg)
	}

	return usedArg, err
}

// --- Option Handlers ---

func handleString(fieldName string) func(*ParameterParser, *OperationConfig, string) error {
	return func(p *ParameterParser, config *OperationConfig, arg string) error {
		switch fieldName {
		case "UserAgent":
			config.UserAgent = arg
		case "CustomRequest":
			config.CustomRequest = arg
		case "UserPassword":
			config.UserPassword = arg
		}
		return nil
	}
}

func handleBool(fieldName string) func(*ParameterParser, *OperationConfig, string) error {
	return func(p *ParameterParser, config *OperationConfig, arg string) error {
		switch fieldName {
		case "InsecureOK":
			config.InsecureOK = true
		case "FollowLocation":
			config.FollowLocation = true
		case "UseHTTPGet":
			config.UseHTTPGet = true
		case "FailOnError":
			config.FailOnError = true
		}
		return nil
	}
}

func handleVerbose(p *ParameterParser, config *OperationConfig, arg string) error {
	return nil
}

func handleHead(p *ParameterParser, config *OperationConfig, arg string) error {
	config.NoBody = true
	config.ShowHeaders = true
	config.UseHTTPGet = false
	return nil
}

func handleURL(p *ParameterParser, config *OperationConfig, arg string) error {
	urlConf := &URLConfig{URL: arg, IsSet: true}
	config.URLList = append(config.URLList, urlConf)
	return nil
}

func handleHeader(p *ParameterParser, config *OperationConfig, arg string) error {
	config.Headers = append(config.Headers, arg)
	return nil
}

func handleData(p *ParameterParser, config *OperationConfig, arg string) error {
	config.PostFields = arg
	return nil
}

func handleOutputFile(p *ParameterParser, config *OperationConfig, arg string) error {
	if len(config.URLList) > 0 {
		config.URLList[len(config.URLList)-1].Outfile = arg
	} else {
		urlConf := &URLConfig{Outfile: arg, IsSet: true}
		config.URLList = append(config.URLList, urlConf)
	}
	return nil
}

func handleRemoteName(p *ParameterParser, config *OperationConfig, arg string) error {
	if len(config.URLList) > 0 {
		config.URLList[len(config.URLList)-1].UseRemote = true
	} else {
		urlConf := &URLConfig{UseRemote: true, IsSet: true}
		config.URLList = append(config.URLList, urlConf)
	}
	return nil
}

func handleConnectTimeout(p *ParameterParser, config *OperationConfig, arg string) error {
	val, err := ParseSecs(arg)
	if err != nil {
		return err
	}
	config.ConnectTimeout = val
	return nil
}

func handleRange(p *ParameterParser, config *OperationConfig, arg string) error {
	if config.UseResume {
		return fmt.Errorf("--continue-at is mutually exclusive with --range")
	}
	// A basic validation, the C code is more complex.
	if !strings.Contains(arg, "-") {
		// curl itself warns but proceeds, we can be stricter.
		return fmt.Errorf("a range must contain at least one dash")
	}
	config.Range = arg
	return nil
}

// parseSizeParameter parses a size string with optional suffix (K, M, G).
func parseSizeParameter(arg string) (int64, error) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return 0, fmt.Errorf("empty size argument")
	}

	i := 0
	for i < len(arg) && unicode.IsDigit(rune(arg[i])) {
		i++
	}
	numPart := arg[:i]
	suffixPart := strings.ToUpper(arg[i:])

	if numPart == "" {
		return 0, fmt.Errorf("invalid number specified for size")
	}

	val, err := strconv.ParseInt(numPart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number for size: %w", err)
	}

	var multiplier int64 = 1
	switch suffixPart {
	case "G":
		multiplier = 1024 * 1024 * 1024
	case "M":
		multiplier = 1024 * 1024
	case "K":
		multiplier = 1024
	case "B", "":
		multiplier = 1
	default:
		return 0, fmt.Errorf("unsupported size unit: %s. Use G, M, K, or B", suffixPart)
	}

	if val > 0 && multiplier > 0 && val > math.MaxInt64/multiplier {
		return 0, fmt.Errorf("size number too large: %s", arg)
	}
	return val * multiplier, nil
}

// parseCertParameter splits a string like "cert:password" into two parts.
func parseCertParameter(param string) (certname, passphrase string) {
	if param == "" {
		return "", ""
	}
	if strings.HasPrefix(param, "pkcs11:") || !strings.ContainsAny(param, ":\\") {
		return param, ""
	}
	var sb strings.Builder
	for i := 0; i < len(param); i++ {
		char := param[i]
		if char == '\\' {
			i++
			if i < len(param) {
				sb.WriteByte(param[i])
			}
		} else if char == ':' {
			if runtime.GOOS == "windows" && sb.Len() == 1 && i+1 < len(param) && (param[i+1] == '\\' || param[i+1] == '/') {
				sb.WriteByte(char)
				continue
			}
			certname = sb.String()
			passphrase = param[i+1:]
			return
		} else {
			sb.WriteByte(char)
		}
	}
	return sb.String(), ""
}
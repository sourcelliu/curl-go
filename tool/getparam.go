package tool

import (
	"fmt"
	"strings"
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
	"url":            {Name: "url", Type: ArgString, Handler: handleURL},
	"verbose":        {Name: "verbose", ShortName: 'v', Type: ArgBool, Handler: handleVerbose},
	"header":         {Name: "header", ShortName: 'H', Type: ArgString, Handler: handleHeader},
	"data":           {Name: "data", ShortName: 'd', Type: ArgString, Handler: handleData},
	"request":        {Name: "request", ShortName: 'X', Type: ArgString, Handler: handleString("CustomRequest")},
	"user-agent":     {Name: "user-agent", ShortName: 'A', Type: ArgString, Handler: handleString("UserAgent")},
	"insecure":       {Name: "insecure", ShortName: 'k', Type: ArgBool, Handler: handleBool("InsecureOK")},
	"location":       {Name: "location", ShortName: 'L', Type: ArgBool, Handler: handleBool("FollowLocation")},
	"output":         {Name: "output", ShortName: 'o', Type: ArgFile, Handler: handleOutputFile},
	"remote-name":    {Name: "remote-name", ShortName: 'O', Type: ArgBool, Handler: handleRemoteName},
	"user":           {Name: "user", ShortName: 'u', Type: ArgString, Handler: handleString("UserPassword")},
	"head":           {Name: "head", ShortName: 'I', Type: ArgBool, Handler: handleHead},
	"get":            {Name: "get", ShortName: 'G', Type: ArgBool, Handler: handleBool("UseHTTPGet")},
}

// shortOptions is a reverse map for finding long options by their short name.
var shortOptions = make(map[rune]Option)

func init() {
	for _, opt := range options {
		if opt.ShortName != 0 {
			shortOptions[opt.ShortName] = opt
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

// ParseOne parses a single flag and its potential argument.
// This is the Go equivalent of the C function `getparameter`.
func (p *ParameterParser) ParseOne(flag, nextarg string) (usedArg bool, err error) {
	var opt Option
	var ok bool
	var arg string

	isLongOpt := strings.HasPrefix(flag, "--")
	if isLongOpt {
		optName := strings.TrimPrefix(flag, "--")
		// TODO: Handle --no- prefix for boolean flags
		opt, ok = options[optName]
		if !ok {
			return false, fmt.Errorf("unknown option: %s", flag)
		}
	} else { // Short option
		shortName := rune(flag[1])
		opt, ok = shortOptions[shortName]
		if !ok {
			return false, fmt.Errorf("unknown option: %s", flag)
		}
		// Check for bundled arguments like -ofoo
		if len(flag) > 2 {
			arg = flag[2:]
		}
	}

	// Check if the option requires an argument
	if opt.Type == ArgString || opt.Type == ArgFile {
		if arg != "" { // Argument was bundled (e.g., -ofoo)
			// The argument is already set
		} else if nextarg != "" {
			arg = nextarg
			usedArg = true
		} else {
			return false, fmt.Errorf("option %s requires an argument", flag)
		}
	}

	// Call the handler for the option
	if opt.Handler != nil {
		err = opt.Handler(p, p.Global.Last, arg)
	}

	return usedArg, err
}

// --- Option Handlers ---

func handleString(fieldName string) func(*ParameterParser, *OperationConfig, string) error {
	return func(p *ParameterParser, config *OperationConfig, arg string) error {
		// This is a simplified handler. A full implementation would use reflection
		// or a map to set the correct field in the OperationConfig.
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
		// For now, we assume arg is empty and we just toggle the boolean to true.
		// A full implementation would handle --no- prefixes.
		switch fieldName {
		case "InsecureOK":
			config.InsecureOK = true
		case "FollowLocation":
			config.FollowLocation = true
		case "UseHTTPGet":
			config.UseHTTPGet = true
		}
		return nil
	}
}

func handleVerbose(p *ParameterParser, config *OperationConfig, arg string) error {
	// For simplicity, we just set a verbose flag. The C code has multiple levels.
	// In our config, this could be a simple boolean or an integer.
	// Let's assume we have a Verbose field.
	// config.Verbose = true
	return nil
}

func handleHead(p *ParameterParser, config *OperationConfig, arg string) error {
	config.NoBody = true
	config.ShowHeaders = true
	config.UseHTTPGet = false // HEAD implies GET
	return nil
}

// The following handlers are placeholders for now.
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
	// This needs to find the right URLConfig to update. For now, assume the last one.
	if len(config.URLList) > 0 {
		config.URLList[len(config.URLList)-1].Outfile = arg
	} else {
		// If no URL is present yet, create one.
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
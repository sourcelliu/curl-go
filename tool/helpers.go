package tool

import "fmt"

// ParameterError is a translation of the C enum `ParameterError` from
// curl-src/src/tool_getparam.h, lines 388-414.
type ParameterError int

const (
	ParamOK ParameterError = iota
	ParamOptionAmbiguous
	ParamOptionUnknown
	ParamRequiresParameter
	ParamBadUse
	ParamHelpRequested
	ParamManualRequested
	ParamVersionInfoRequested
	ParamEnginesRequested
	ParamCAEmbedRequested
	ParamGotExtraParameter
	ParamBadNumeric
	ParamNegativeNumeric
	ParamLibcurlDoesntSupport
	ParamLibcurlUnsupportedProtocol
	ParamNoMem
	ParamNextOperation
	ParamNoPrefix
	ParamNumberTooLarge
	ParamContDispResumeFrom
	ParamReadError
	ParamExpandError
	ParamBlankString
	ParamVarSyntax
	ParamLast // Keep last for counting
)

// String provides a human-readable representation of a ParameterError.
// This is the idiomatic Go replacement for the C function `param2text` from
// curl-src/src/tool_helpers.c, lines 35-71.
func (e ParameterError) String() string {
	switch e {
	case ParamOK:
		return "ok"
	case ParamGotExtraParameter:
		return "had unsupported trailing garbage"
	case ParamOptionUnknown:
		return "is unknown"
	case ParamOptionAmbiguous:
		return "is ambiguous"
	case ParamRequiresParameter:
		return "requires parameter"
	case ParamBadUse:
		return "is badly used here"
	case ParamBadNumeric:
		return "expected a proper numerical parameter"
	case ParamNegativeNumeric:
		return "expected a positive numerical parameter"
	case ParamLibcurlDoesntSupport:
		return "the installed libcurl version does not support this"
	case ParamLibcurlUnsupportedProtocol:
		return "a specified protocol is unsupported by libcurl"
	case ParamNoMem:
		return "out of memory"
	case ParamNoPrefix:
		return "the given option cannot be reversed with a --no- prefix"
	case ParamNumberTooLarge:
		return "too large number"
	case ParamContDispResumeFrom:
		return "--continue-at and --remote-header-name cannot be combined"
	case ParamReadError:
		return "error encountered when reading a file"
	case ParamExpandError:
		return "variable expansion failure"
	case ParamBlankString:
		return "blank argument where content is expected"
	case ParamVarSyntax:
		return "syntax error in --variable argument"
	default:
		return "unknown error"
	}
}

// HTTPRequest is a translation of the C enum `HttpReq` from
// curl-src/src/tool_sdecls.h, lines 112-119.
type HTTPRequest int

const (
	HTTPRequestUnspec HTTPRequest = iota
	HTTPRequestGet
	HTTPRequestHead
	HTTPRequestMimePost
	HTTPRequestSimplePost
	HTTPRequestPut
)

// String provides a human-readable representation of an HTTPRequest.
func (r HTTPRequest) String() string {
	switch r {
	case HTTPRequestGet:
		return "GET (-G, --get)"
	case HTTPRequestHead:
		return "HEAD (-I, --head)"
	case HTTPRequestMimePost:
		return "multipart formpost (-F, --form)"
	case HTTPRequestSimplePost:
		return "POST (-d, --data)"
	case HTTPRequestPut:
		return "PUT (-T, --upload-file)"
	default:
		return ""
	}
}

// HTTPRequestManager helps manage the HTTP request type state.
// This is the idiomatic Go replacement for the C function `SetHTTPrequest`
// from curl-src/src/tool_helpers.c, lines 73-93.
type HTTPRequestManager struct {
	request HTTPRequest
}

// Set sets the HTTP request type. It returns an error if a conflicting
// request type has already been set.
func (m *HTTPRequestManager) Set(req HTTPRequest) error {
	if m.request == HTTPRequestUnspec || m.request == req {
		m.request = req
		return nil
	}
	return fmt.Errorf("you can only select one HTTP request method! You asked for both %s and %s",
		req.String(), m.request.String())
}

// CustomRequestHelper provides helpful warnings for the use of custom requests.
// This is a translation of the C function `customrequest_helper` from
// curl-src/src/tool_helpers.c, lines 95-119.
func CustomRequestHelper(req HTTPRequest, method string) string {
	defaultMethods := map[HTTPRequest]string{
		HTTPRequestGet:        "GET",
		HTTPRequestHead:       "HEAD",
		HTTPRequestMimePost:   "POST",
		HTTPRequestSimplePost: "POST",
		HTTPRequestPut:        "PUT",
	}

	if method == "" {
		return ""
	}

	if defaultMethod, ok := defaultMethods[req]; ok && method == defaultMethod {
		return fmt.Sprintf("Unnecessary use of -X or --request, %s is already inferred.", defaultMethod)
	}

	if method == "HEAD" {
		return "Setting custom HTTP method to HEAD with -X/--request may not work the way you want. Consider using -I/--head instead."
	}

	return ""
}
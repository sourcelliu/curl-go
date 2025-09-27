package tool

// This file contains the Go translation of the central `OperationConfig`
// struct from `curl-src/src/tool_cfgable.h`.

// URLConfig holds the configuration for a single URL to be fetched.
// It is a translation of the C `getout` struct.
type URLConfig struct {
	URL      string
	Outfile  string
	Infile   string
	IsSet    bool // Tracks if this node has been used
	NoGlob   bool
	UseRemote bool
}

// OperationConfig holds all the settings for a single curl operation.
// It is the Go equivalent of the C `OperationConfig` struct.
// For now, it contains a subset of the fields to support initial parsing.
type OperationConfig struct {
	// String options
	UserAgent         string
	CookieJar         string
	PostFields        string
	Referer           string
	UserPassword      string
	ProxyUserPassword string
	Proxy             string
	HeaderFile        string
	WriteOut          string
	Range             string
	CustomRequest     string

	// Slices of strings
	Headers     []string

	// Numeric options
	MaxRedirs      int64
	AuthType       uint // Bitmask
	FollowLocation bool

	// Boolean options
	InsecureOK         bool
	ShowHeaders        bool
	NoBody             bool
	UseHTTPGet         bool
	ContentDisposition bool

	// URL List
	URLList []*URLConfig

	// Linked list for multiple operations
	Next *OperationConfig
	Prev *OperationConfig
}

// NewOperationConfig creates and returns a new, initialized OperationConfig.
func NewOperationConfig() *OperationConfig {
	return &OperationConfig{
		URLList: make([]*URLConfig, 0),
	}
}

// GlobalConfig holds settings that apply to all operations.
type GlobalConfig struct {
	First *OperationConfig
	Last  *OperationConfig
}

// NewGlobalConfig creates a new GlobalConfig, initializes it, and sets up
// the first OperationConfig.
func NewGlobalConfig() *GlobalConfig {
	g := &GlobalConfig{}
	first := NewOperationConfig()
	g.First = first
	g.Last = first
	return g
}
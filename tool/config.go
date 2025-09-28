package tool

import "time"

// This file contains the Go translation of the central `OperationConfig`
// and `GlobalConfig` structs from `curl-src/src/tool_cfgable.h`.
//
// It also serves as the Go equivalent for the C file `tool_cfgable.c`,
// which is responsible for the allocation and deallocation of these structs.
// In Go, allocation is handled by constructor functions (`New...`), and
// deallocation is handled automatically by the garbage collector, making
// explicit `free` functions unnecessary.

// URLConfig holds the configuration for a single URL to be fetched.
// It is a translation of the C `getout` struct.
type URLConfig struct {
	URL       string
	Outfile   string
	Infile    string
	IsSet     bool // Tracks if this node has been used
	NoGlob    bool
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
	Headers []string

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
	RemoteTime         bool
	FailOnError        bool
	UseResume          bool

	// Timeouts
	ConnectTimeout time.Duration

	// URL List
	URLList []*URLConfig

	// Linked list for multiple operations
	Next *OperationConfig
	Prev *OperationConfig
}

// NewOperationConfig creates and returns a new, initialized OperationConfig.
// This is the Go equivalent of the C function `config_alloc`.
func NewOperationConfig() *OperationConfig {
	return &OperationConfig{
		// Initialize fields with their default zero values, which is often correct.
		// Specific defaults can be set here if needed.
		URLList: make([]*URLConfig, 0),
	}
}

// GlobalConfig holds settings that apply to all operations.
// It is the Go equivalent of the C `GlobalConfig` struct.
type GlobalConfig struct {
	First *OperationConfig
	Last  *OperationConfig
	// Other global fields like TraceDump, LibCurl, etc., will be added here as needed.
}

// NewGlobalConfig creates a new GlobalConfig, initializes it, and sets up
// the first OperationConfig. This is the Go equivalent of `globalconf_init`.
func NewGlobalConfig() *GlobalConfig {
	g := &GlobalConfig{}
	first := NewOperationConfig()
	g.First = first
	g.Last = first
	return g
}

// Note: The C file `tool_cfgable.c` contains `config_free` and
// `free_config_fields`. These are not needed in Go because the garbage
// collector automatically handles deallocation when the structs are no longer
// referenced.
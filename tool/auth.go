package tool

// AuthType is a bitmask for authentication methods.
// It is a translation of the CURLAUTH_* defines used in curl.
type AuthType uint

const (
	AuthNone     AuthType = 0
	AuthBasic    AuthType = 1 << iota
	AuthDigest            // 1 << 1
	AuthNegotiate         // 1 << 2
	AuthNTLM              // 1 << 3
	// ... other auth types can be added here
	AuthAny = ^AuthType(0) // Represents any authentication method
)
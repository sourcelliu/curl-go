package tool

import (
	"runtime"
	"testing"
)

func TestGetInfo(t *testing.T) {
	t.Run("first call populates data", func(t *testing.T) {
		info := GetInfo()
		if info == nil {
			t.Fatal("GetInfo() returned nil")
		}

		if info.GoVersion != runtime.Version() {
			t.Errorf("GoVersion = %q; want %q", info.GoVersion, runtime.Version())
		}
		if info.OS != runtime.GOOS {
			t.Errorf("OS = %q; want %q", info.OS, runtime.GOOS)
		}
		if info.Arch != runtime.GOARCH {
			t.Errorf("Arch = %q; want %q", info.Arch, runtime.GOARCH)
		}

		// Check for a required protocol
		foundHttp := false
		for _, p := range info.Protocols {
			if p == "http" {
				foundHttp = true
				break
			}
		}
		if !foundHttp {
			t.Error("Protocols slice should contain 'http'")
		}

		// Check for a required feature
		if ssl, ok := info.Features["SSL"]; !ok || !ssl {
			t.Error("Features map should indicate SSL is enabled")
		}
	})

	t.Run("is a singleton", func(t *testing.T) {
		info1 := GetInfo()
		info2 := GetInfo()

		// The pointers should be identical, proving that the struct was
		// only initialized once.
		if info1 != info2 {
			t.Error("GetInfo() should return the same instance on subsequent calls")
		}
	})
}
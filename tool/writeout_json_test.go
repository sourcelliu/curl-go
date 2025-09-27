package tool

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"runtime"
	"testing"
)

func TestWriteOutJSON(t *testing.T) {
	// Sample data similar to what would be collected after a curl transfer.
	sampleData := map[string]interface{}{
		"url_effective": "http://example.com",
		"http_code":     200,
		"time_total":    1.23,
		"size_download": 1024,
	}

	var buf bytes.Buffer
	if err := WriteOutJSON(&buf, sampleData); err != nil {
		t.Fatalf("WriteOutJSON() returned an unexpected error: %v", err)
	}

	// Unmarshal the output back into a map to verify its contents.
	// This is more robust than string comparison, as key order is not guaranteed.
	var resultMap map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &resultMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON output: %v\nOutput was: %s", err, buf.String())
	}

	// Verify that the special "curl_version" key was added.
	if version, ok := resultMap["curl_version"]; !ok || version != runtime.Version() {
		t.Errorf("Expected 'curl_version' key to be %q, but it was %q", runtime.Version(), version)
	}
	// Delete it so we can compare the rest of the map.
	delete(resultMap, "curl_version")

	// The numeric types might get unmarshaled as float64, so we need to
	// convert our original data for comparison.
	expectedMap := make(map[string]interface{})
	for k, v := range sampleData {
		switch v := v.(type) {
		case int:
			expectedMap[k] = float64(v)
		default:
			expectedMap[k] = v
		}
	}

	if !reflect.DeepEqual(resultMap, expectedMap) {
		t.Errorf("Result map differs from expected map.\nGot:    %v\nWanted: %v", resultMap, expectedMap)
	}
}

func TestHeaderJSON(t *testing.T) {
	// Sample http.Header object.
	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"X-Request-Id": []string{"12345"},
		"Set-Cookie":   []string{"a=1", "b=2"}, // Header with multiple values
	}

	var buf bytes.Buffer
	if err := HeaderJSON(&buf, headers); err != nil {
		t.Fatalf("HeaderJSON() returned an unexpected error: %v", err)
	}

	// Unmarshal the output and verify its contents.
	var resultMap http.Header
	if err := json.Unmarshal(buf.Bytes(), &resultMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON output: %v\nOutput was: %s", err, buf.String())
	}

	if !reflect.DeepEqual(resultMap, headers) {
		t.Errorf("Result map differs from expected map.\nGot:    %v\nWanted: %v", resultMap, headers)
	}
}
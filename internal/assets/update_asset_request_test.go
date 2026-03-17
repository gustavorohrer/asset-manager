package assets

import (
	"strings"
	"testing"
	"time"
)

func TestParseUpdateAssetRequestBodySuccess(t *testing.T) {
	body := []byte(`{"name":"  Updated name  ","description":"updated description","lastScan":"2024-10-08T00:00:00Z"}`)

	input, details := ParseUpdateAssetRequestBody(body)
	if len(details) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(details))
	}
	if input.Name == nil || *input.Name != "Updated name" {
		t.Fatalf("expected trimmed name, got %#v", input.Name)
	}
	if input.Description == nil || *input.Description != "updated description" {
		t.Fatalf("expected description, got %#v", input.Description)
	}
	if !input.LastScanSet || input.LastScan == nil {
		t.Fatalf("expected lastScan set, got LastScanSet=%v LastScan=%v", input.LastScanSet, input.LastScan)
	}
	if !input.LastScan.Equal(time.Date(2024, 10, 8, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected lastScan value: %v", input.LastScan)
	}
}

func TestParseUpdateAssetRequestBodyAllowsNullLastScan(t *testing.T) {
	body := []byte(`{"lastScan":null}`)

	input, details := ParseUpdateAssetRequestBody(body)
	if len(details) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(details))
	}
	if !input.LastScanSet {
		t.Fatal("expected LastScanSet=true")
	}
	if input.LastScan != nil {
		t.Fatalf("expected nil lastScan, got %v", input.LastScan)
	}
}

func TestParseUpdateAssetRequestBodyRejectsUnknownField(t *testing.T) {
	_, details := ParseUpdateAssetRequestBody([]byte(`{"foo":"bar"}`))
	if len(details) == 0 {
		t.Fatal("expected validation errors")
	}
}

func TestParseUpdateAssetRequestBodyRejectsNullDescription(t *testing.T) {
	_, details := ParseUpdateAssetRequestBody([]byte(`{"description":null}`))
	if len(details) == 0 {
		t.Fatal("expected validation errors")
	}
}

func TestParseUpdateAssetRequestBodyRejectsInvalidJSON(t *testing.T) {
	_, details := ParseUpdateAssetRequestBody([]byte(`{`))
	if len(details) == 0 {
		t.Fatal("expected validation errors")
	}
}

func TestParseUpdateAssetRequestBodyRejectsEmptyName(t *testing.T) {
	_, details := ParseUpdateAssetRequestBody([]byte(`{"name":"   "}`))
	if len(details) == 0 {
		t.Fatal("expected validation errors")
	}
}

func TestParseUpdateAssetRequestBodyRejectsLongDescription(t *testing.T) {
	value := strings.Repeat("a", 10001)
	body := `{"description":"` + value + `"}`

	_, details := ParseUpdateAssetRequestBody([]byte(body))
	if len(details) == 0 {
		t.Fatal("expected validation errors")
	}
}

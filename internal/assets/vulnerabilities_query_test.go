package assets

import (
	"net/url"
	"testing"
)

func TestParseListAssetVulnerabilitiesQueryDefaults(t *testing.T) {
	query, details := ParseListAssetVulnerabilitiesQuery(url.Values{})
	if len(details) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(details))
	}
	if query.Page != 1 {
		t.Fatalf("expected page 1, got %d", query.Page)
	}
	if query.PageSize != 20 {
		t.Fatalf("expected pageSize 20, got %d", query.PageSize)
	}
	if query.Severity != nil {
		t.Fatalf("expected nil severity, got %v", *query.Severity)
	}
}

func TestParseListAssetVulnerabilitiesQueryValidSeverity(t *testing.T) {
	query, details := ParseListAssetVulnerabilitiesQuery(url.Values{
		"severity": {"critical"},
		"page":     {"2"},
		"pageSize": {"10"},
	})

	if len(details) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(details))
	}
	if query.Severity == nil || *query.Severity != SeverityCritical {
		t.Fatalf("expected severity CRITICAL, got %v", query.Severity)
	}
	if query.Page != 2 {
		t.Fatalf("expected page 2, got %d", query.Page)
	}
	if query.PageSize != 10 {
		t.Fatalf("expected pageSize 10, got %d", query.PageSize)
	}
}

func TestParseListAssetVulnerabilitiesQueryInvalid(t *testing.T) {
	_, details := ParseListAssetVulnerabilitiesQuery(url.Values{
		"page":     {"0"},
		"pageSize": {"101"},
		"severity": {"urgent"},
	})

	if len(details) != 3 {
		t.Fatalf("expected 3 validation errors, got %d", len(details))
	}
}

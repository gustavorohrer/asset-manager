package assets

import (
	"net/url"
	"testing"
	"time"
)

func TestParseListAssetsQueryDefaults(t *testing.T) {
	query, details := ParseListAssetsQuery(url.Values{})
	if len(details) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(details))
	}

	if query.Page != 1 {
		t.Fatalf("expected default page 1, got %d", query.Page)
	}
	if query.PageSize != 20 {
		t.Fatalf("expected default pageSize 20, got %d", query.PageSize)
	}
	if query.SortBy != SortByCreatedAt {
		t.Fatalf("expected default sortBy createdAt, got %s", query.SortBy)
	}
	if query.SortOrder != SortOrderDesc {
		t.Fatalf("expected default sortOrder desc, got %s", query.SortOrder)
	}
	if query.HasVulnerabilities != nil {
		t.Fatalf("expected has_vulnerabilities nil by default, got %v", *query.HasVulnerabilities)
	}
	if query.HasThreats != nil {
		t.Fatalf("expected has_threats nil by default, got %v", *query.HasThreats)
	}
}

func TestParseListAssetsQueryValidValues(t *testing.T) {
	values := url.Values{
		"name":                {"router"},
		"created_from":        {"2024-01-01T00:00:00Z"},
		"created_to":          {"2024-12-31T00:00:00Z"},
		"last_scan_from":      {"2024-05-01T00:00:00Z"},
		"last_scan_to":        {"2024-11-01T00:00:00Z"},
		"page":                {"2"},
		"pageSize":            {"50"},
		"sortBy":              {"name"},
		"sortOrder":           {"asc"},
		"has_vulnerabilities": {"true"},
		"has_threats":         {"false"},
		"unknown":             {"ignored"},
	}

	query, details := ParseListAssetsQuery(values)
	if len(details) != 0 {
		t.Fatalf("expected no validation errors, got %d: %#v", len(details), details)
	}

	if query.NameContains != "router" {
		t.Fatalf("expected name filter router, got %q", query.NameContains)
	}
	if query.Page != 2 {
		t.Fatalf("expected page 2, got %d", query.Page)
	}
	if query.PageSize != 50 {
		t.Fatalf("expected pageSize 50, got %d", query.PageSize)
	}
	if query.SortBy != SortByName {
		t.Fatalf("expected sortBy name, got %s", query.SortBy)
	}
	if query.SortOrder != SortOrderAsc {
		t.Fatalf("expected sortOrder asc, got %s", query.SortOrder)
	}
	if query.CreatedFrom == nil || query.CreatedTo == nil || query.LastScanFrom == nil || query.LastScanTo == nil {
		t.Fatal("expected all date filters to be parsed")
	}
	if query.HasVulnerabilities == nil || !*query.HasVulnerabilities {
		t.Fatalf("expected has_vulnerabilities=true, got %#v", query.HasVulnerabilities)
	}
	if query.HasThreats == nil || *query.HasThreats {
		t.Fatalf("expected has_threats=false, got %#v", query.HasThreats)
	}
}

func TestParseListAssetsQueryInvalidValues(t *testing.T) {
	values := url.Values{
		"page":                {"0"},
		"pageSize":            {"101"},
		"sortBy":              {"invalid"},
		"sortOrder":           {"up"},
		"created_from":        {"bad-date"},
		"has_vulnerabilities": {"truthy"},
		"has_threats":         {"falsy"},
	}

	_, details := ParseListAssetsQuery(values)
	if len(details) == 0 {
		t.Fatal("expected validation errors")
	}

	hasField := func(field string) bool {
		for _, d := range details {
			if d.Field == field {
				return true
			}
		}
		return false
	}

	for _, expectedField := range []string{"page", "pageSize", "sortBy", "sortOrder", "created_from", "has_vulnerabilities", "has_threats"} {
		if !hasField(expectedField) {
			t.Fatalf("expected validation error for field %s", expectedField)
		}
	}
}

func TestParseListAssetsQueryPageTooLarge(t *testing.T) {
	_, details := ParseListAssetsQuery(url.Values{
		"page": {"10001"},
	})

	if len(details) != 1 {
		t.Fatalf("expected 1 validation error, got %d", len(details))
	}
	if details[0].Field != "page" {
		t.Fatalf("expected page validation error, got %s", details[0].Field)
	}
}

func TestParseListAssetsQueryInvalidRanges(t *testing.T) {
	createdFrom := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	createdTo := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	_, details := ParseListAssetsQuery(url.Values{
		"created_from": {createdFrom},
		"created_to":   {createdTo},
	})

	if len(details) != 1 {
		t.Fatalf("expected 1 validation error, got %d", len(details))
	}
	if details[0].Field != "created_from" {
		t.Fatalf("expected created_from validation error, got %s", details[0].Field)
	}
}

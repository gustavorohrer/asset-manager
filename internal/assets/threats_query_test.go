package assets

import (
	"net/url"
	"testing"
)

func TestParseListAssetThreatsQueryDefaults(t *testing.T) {
	query, details := ParseListAssetThreatsQuery(url.Values{})
	if len(details) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(details))
	}
	if query.Page != 1 {
		t.Fatalf("expected page 1, got %d", query.Page)
	}
	if query.PageSize != 20 {
		t.Fatalf("expected pageSize 20, got %d", query.PageSize)
	}
	if query.RiskLevel != nil {
		t.Fatalf("expected nil riskLevel, got %v", *query.RiskLevel)
	}
}

func TestParseListAssetThreatsQueryValidRiskLevel(t *testing.T) {
	query, details := ParseListAssetThreatsQuery(url.Values{
		"riskLevel": {"medium"},
		"page":      {"2"},
		"pageSize":  {"10"},
	})

	if len(details) != 0 {
		t.Fatalf("expected no validation errors, got %d", len(details))
	}
	if query.RiskLevel == nil || *query.RiskLevel != RiskLevelMedium {
		t.Fatalf("expected riskLevel MEDIUM, got %v", query.RiskLevel)
	}
	if query.Page != 2 {
		t.Fatalf("expected page 2, got %d", query.Page)
	}
	if query.PageSize != 10 {
		t.Fatalf("expected pageSize 10, got %d", query.PageSize)
	}
}

func TestParseListAssetThreatsQueryInvalid(t *testing.T) {
	_, details := ParseListAssetThreatsQuery(url.Values{
		"page":      {"0"},
		"pageSize":  {"101"},
		"riskLevel": {"critical"},
	})

	if len(details) != 3 {
		t.Fatalf("expected 3 validation errors, got %d", len(details))
	}
}

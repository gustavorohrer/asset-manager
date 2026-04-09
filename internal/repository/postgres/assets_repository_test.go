package postgres

import (
	"strings"
	"testing"
	"time"

	"github.com/gustavorohrer/ecl-be-challenge/internal/assets"
)

func TestEscapeLikeLiteral(t *testing.T) {
	got := escapeLikeLiteral(`a\b%c_d`)
	want := `a\\b\%c\_d`
	if got != want {
		t.Fatalf("unexpected escaped value. want=%q got=%q", want, got)
	}
}

func TestDateForSQL(t *testing.T) {
	value := time.Date(2024, 10, 8, 23, 59, 59, 0, time.FixedZone("-0300", -3*60*60))
	got := dateForSQL(value)
	if got != "2024-10-08" {
		t.Fatalf("unexpected SQL date. want=2024-10-08 got=%s", got)
	}
}

func TestBuildFiltersNameEscapingAndDates(t *testing.T) {
	createdFrom := time.Date(2024, 1, 1, 13, 0, 0, 0, time.FixedZone("+0200", 2*60*60))
	createdTo := time.Date(2024, 12, 31, 23, 0, 0, 0, time.UTC)

	whereClause, args := buildFilters(assets.ListAssetsQuery{
		NameContains: `router_%\v2`,
		CreatedFrom:  &createdFrom,
		CreatedTo:    &createdTo,
	})

	if !strings.Contains(whereClause, "a.name ILIKE $1 ESCAPE E'\\\\'") {
		t.Fatalf("expected escaped ILIKE condition, got: %s", whereClause)
	}
	if !strings.Contains(whereClause, "a.createdat >= $2::date") {
		t.Fatalf("expected created_from condition, got: %s", whereClause)
	}
	if !strings.Contains(whereClause, "a.createdat <= $3::date") {
		t.Fatalf("expected created_to condition, got: %s", whereClause)
	}

	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}

	if args[0] != `%router\_\%\\v2%` {
		t.Fatalf("unexpected escaped LIKE arg: %q", args[0])
	}
	if args[1] != "2024-01-01" {
		t.Fatalf("unexpected created_from arg: %v", args[1])
	}
	if args[2] != "2024-12-31" {
		t.Fatalf("unexpected created_to arg: %v", args[2])
	}
}

func TestBuildOrder(t *testing.T) {
	got := buildOrder("fa", assets.SortByLastScan, assets.SortOrderAsc)
	want := "fa.lastscan ASC NULLS LAST, fa.id ASC"
	if got != want {
		t.Fatalf("unexpected order clause. want=%q got=%q", want, got)
	}
}

func TestBuildFiltersHasFindingsFlags(t *testing.T) {
	hasVulnerabilities := true
	hasThreats := false

	whereClause, args := buildFilters(assets.ListAssetsQuery{
		HasVulnerabilities: &hasVulnerabilities,
		HasThreats:         &hasThreats,
	})

	if !strings.Contains(whereClause, "JOIN vulnerability v ON v.scanid = latest_component_scans.scanid") {
		t.Fatalf("expected vulnerabilities filter condition, got: %s", whereClause)
	}
	if !strings.Contains(whereClause, "JOIN threat t ON t.scanid = latest_component_scans.scanid") {
		t.Fatalf("expected threats filter condition, got: %s", whereClause)
	}
	if !strings.Contains(whereClause, ") = $1") {
		t.Fatalf("expected has_vulnerabilities placeholder, got: %s", whereClause)
	}
	if !strings.Contains(whereClause, ") = $2") {
		t.Fatalf("expected has_threats placeholder, got: %s", whereClause)
	}

	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}
	if args[0] != true {
		t.Fatalf("expected first arg true, got %v", args[0])
	}
	if args[1] != false {
		t.Fatalf("expected second arg false, got %v", args[1])
	}
}

func TestBuildFiltersHasFindingsCombinedFlag(t *testing.T) {
	hasFindings := true

	whereClause, args := buildFilters(assets.ListAssetsQuery{
		HasFindings: &hasFindings,
	})

	if !strings.Contains(whereClause, "JOIN vulnerability v ON v.scanid = latest_component_scans.scanid") {
		t.Fatalf("expected vulnerabilities in combined findings condition, got: %s", whereClause)
	}
	if !strings.Contains(whereClause, "JOIN threat t ON t.scanid = latest_component_scans.scanid") {
		t.Fatalf("expected threats in combined findings condition, got: %s", whereClause)
	}
	if !strings.Contains(whereClause, "OR") {
		t.Fatalf("expected OR in combined findings condition, got: %s", whereClause)
	}
	if !strings.Contains(whereClause, ") = $1") {
		t.Fatalf("expected has_findings placeholder, got: %s", whereClause)
	}

	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args))
	}
	if args[0] != true {
		t.Fatalf("expected arg true, got %v", args[0])
	}
}

func TestBuildListAssetsDataSQLUsesSeparatedAggregations(t *testing.T) {
	sql := buildListAssetsDataSQL(
		"WHERE a.name ILIKE $1",
		"fa.createdat DESC, fa.id ASC",
		"$2",
		"$3",
	)

	if !strings.Contains(sql, "vulnerability_counts_by_asset AS") {
		t.Fatalf("expected vulnerability_counts_by_asset CTE, got: %s", sql)
	}
	if !strings.Contains(sql, "threat_presence_by_asset AS") {
		t.Fatalf("expected threat_presence_by_asset CTE, got: %s", sql)
	}
	if !strings.Contains(sql, "COUNT(*) FILTER (WHERE v.severity = 'HIGH')") {
		t.Fatalf("expected HIGH vulnerability counter, got: %s", sql)
	}
	if !strings.Contains(sql, "COUNT(*) FILTER (WHERE v.severity = 'MEDIUM')") {
		t.Fatalf("expected MEDIUM vulnerability counter, got: %s", sql)
	}
	if !strings.Contains(sql, "COALESCE(vca.vulnerabilities_total > 0, FALSE) AS has_vulnerabilities") {
		t.Fatalf("expected has_vulnerabilities derived from vulnerability counts, got: %s", sql)
	}
	if !strings.Contains(sql, "COALESCE(vca.vulnerabilities_high, 0) AS vulnerabilities_high") {
		t.Fatalf("expected vulnerabilities_high projection, got: %s", sql)
	}
	if !strings.Contains(sql, "COALESCE(vca.vulnerabilities_medium, 0) AS vulnerabilities_medium") {
		t.Fatalf("expected vulnerabilities_medium projection, got: %s", sql)
	}
	if !strings.Contains(sql, "COALESCE(vca.vulnerabilities_total, 0) AS vulnerabilities_total") {
		t.Fatalf("expected vulnerabilities_total projection, got: %s", sql)
	}
}

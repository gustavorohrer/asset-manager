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

package main

import "testing"

func TestResolveDatabaseURLFromDATABASEURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db?sslmode=disable")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")

	got, err := resolveDatabaseURL()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "postgres://user:pass@localhost:5432/db?sslmode=disable" {
		t.Fatalf("unexpected database URL: %s", got)
	}
}

func TestResolveDatabaseURLFromDBVars(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "eclypsiumdb")
	t.Setenv("DB_USER", "applicant")
	t.Setenv("DB_PASSWORD", "goodluck")

	got, err := resolveDatabaseURL()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "postgres://applicant:goodluck@localhost:5432/eclypsiumdb?sslmode=disable"
	if got != want {
		t.Fatalf("unexpected database URL. want=%s got=%s", want, got)
	}
}

func TestResolveDatabaseURLMissingConfig(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")

	_, err := resolveDatabaseURL()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

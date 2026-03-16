package domain

import (
	"fmt"
	"strings"
	"time"
)

type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

func (s Severity) IsValid() bool {
	switch s {
	case SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical:
		return true
	default:
		return false
	}
}

func ParseSeverity(value string) (Severity, error) {
	s := Severity(strings.ToUpper(strings.TrimSpace(value)))
	if !s.IsValid() {
		return "", fmt.Errorf("invalid severity: %q", value)
	}
	return s, nil
}

type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "LOW"
	RiskLevelMedium RiskLevel = "MEDIUM"
	RiskLevelHigh   RiskLevel = "HIGH"
)

func (r RiskLevel) IsValid() bool {
	switch r {
	case RiskLevelLow, RiskLevelMedium, RiskLevelHigh:
		return true
	default:
		return false
	}
}

func ParseRiskLevel(value string) (RiskLevel, error) {
	r := RiskLevel(strings.ToUpper(strings.TrimSpace(value)))
	if !r.IsValid() {
		return "", fmt.Errorf("invalid risk level: %q", value)
	}
	return r, nil
}

type Asset struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	LastScan    *time.Time
	Components  []Component
}

type Component struct {
	ID        string
	Name      string
	Version   string
	Vendor    string
	Type      string
	CreatedAt time.Time
	LastScan  *time.Time
	AssetID   string
	Scans     []Scan
}

type Scan struct {
	ID              string
	PerformedAt     time.Time
	ScannerName     string
	ComponentID     string
	Vulnerabilities []Vulnerability
	Threats         []Threat
}

type Vulnerability struct {
	ID          string
	Description string
	Severity    Severity
	ScanID      string
}

type Threat struct {
	ID          string
	Description string
	RiskLevel   RiskLevel
	Type        string
	ScanID      string
}

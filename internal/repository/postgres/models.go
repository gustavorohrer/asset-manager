package postgres

import (
	"fmt"
	"time"

	"github.com/gustavorohrer/ecl-be-challenge/internal/domain"
)

type assetRow struct {
	ID          string     `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"createdat"`
	LastScan    *time.Time `db:"lastscan"`
}

func (r assetRow) toDomain() domain.Asset {
	return domain.Asset{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		LastScan:    r.LastScan,
	}
}

type componentRow struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Version   string     `db:"version"`
	Vendor    string     `db:"vendor"`
	Type      string     `db:"type"`
	CreatedAt time.Time  `db:"createdat"`
	LastScan  *time.Time `db:"lastscan"`
	AssetID   string     `db:"assetid"`
}

func (r componentRow) toDomain() domain.Component {
	return domain.Component{
		ID:        r.ID,
		Name:      r.Name,
		Version:   r.Version,
		Vendor:    r.Vendor,
		Type:      r.Type,
		CreatedAt: r.CreatedAt,
		LastScan:  r.LastScan,
		AssetID:   r.AssetID,
	}
}

type scanRow struct {
	ID          string    `db:"id"`
	PerformedAt time.Time `db:"performedat"`
	ScannerName string    `db:"scannername"`
	ComponentID string    `db:"componentid"`
}

func (r scanRow) toDomain() domain.Scan {
	return domain.Scan{
		ID:          r.ID,
		PerformedAt: r.PerformedAt,
		ScannerName: r.ScannerName,
		ComponentID: r.ComponentID,
	}
}

type vulnerabilityRow struct {
	ID          string `db:"id"`
	Description string `db:"description"`
	Severity    string `db:"severity"`
	ScanID      string `db:"scanid"`
}

func (r vulnerabilityRow) toDomain() (domain.Vulnerability, error) {
	severity, err := domain.ParseSeverity(r.Severity)
	if err != nil {
		return domain.Vulnerability{}, fmt.Errorf("parse vulnerability severity: %w", err)
	}

	return domain.Vulnerability{
		ID:          r.ID,
		Description: r.Description,
		Severity:    severity,
		ScanID:      r.ScanID,
	}, nil
}

type threatRow struct {
	ID          string `db:"id"`
	Description string `db:"description"`
	RiskLevel   string `db:"risklevel"`
	Type        string `db:"type"`
	ScanID      string `db:"scanid"`
}

func (r threatRow) toDomain() (domain.Threat, error) {
	riskLevel, err := domain.ParseRiskLevel(r.RiskLevel)
	if err != nil {
		return domain.Threat{}, fmt.Errorf("parse threat risk level: %w", err)
	}

	return domain.Threat{
		ID:          r.ID,
		Description: r.Description,
		RiskLevel:   riskLevel,
		Type:        r.Type,
		ScanID:      r.ScanID,
	}, nil
}

package assets

import (
	"errors"
	"time"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
	maxPage         = 10000
)

type SortBy string

const (
	SortByCreatedAt SortBy = "createdAt"
	SortByName      SortBy = "name"
	SortByLastScan  SortBy = "lastScan"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "LOW"
	RiskLevelMedium RiskLevel = "MEDIUM"
	RiskLevelHigh   RiskLevel = "HIGH"
)

type ListAssetsQuery struct {
	NameContains string

	CreatedFrom  *time.Time
	CreatedTo    *time.Time
	LastScanFrom *time.Time
	LastScanTo   *time.Time

	HasVulnerabilities *bool
	HasThreats         *bool
	HasFindings        *bool

	Page      int
	PageSize  int
	SortBy    SortBy
	SortOrder SortOrder
}

type AssetSummary struct {
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	CreatedAt           time.Time           `json:"createdAt"`
	LastScan            *time.Time          `json:"lastScan"`
	HasVulnerabilities  bool                `json:"hasVulnerabilities"`
	HasThreats          bool                `json:"hasThreats"`
	VulnerabilityCounts VulnerabilityCounts `json:"vulnerabilityCounts"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type ListAssetsResponse struct {
	Data       []AssetSummary `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

type VulnerabilityCounts struct {
	High   int `json:"high"`
	Medium int `json:"medium"`
	Total  int `json:"total"`
}

type AssetRiskSummary struct {
	Total               int `json:"total"`
	WithVulnerabilities int `json:"withVulnerabilities"`
	WithThreats         int `json:"withThreats"`
}

type AssetComponent struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Version   string     `json:"version"`
	Vendor    string     `json:"vendor"`
	Type      string     `json:"type"`
	CreatedAt time.Time  `json:"createdAt"`
	LastScan  *time.Time `json:"lastScan"`
	AssetID   string     `json:"assetId"`
}

type AssetDetails struct {
	ID                 string           `json:"id"`
	Name               string           `json:"name"`
	Description        string           `json:"description"`
	CreatedAt          time.Time        `json:"createdAt"`
	LastScan           *time.Time       `json:"lastScan"`
	HasVulnerabilities bool             `json:"hasVulnerabilities"`
	HasThreats         bool             `json:"hasThreats"`
	Components         []AssetComponent `json:"components"`
}

type AssetUpdated struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	LastScan    *time.Time `json:"lastScan"`
}

type AssetDeleted struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type ListAssetVulnerabilitiesQuery struct {
	Page     int
	PageSize int
	Severity *Severity
}

type ListAssetThreatsQuery struct {
	Page      int
	PageSize  int
	RiskLevel *RiskLevel
}

type UpdateAssetInput struct {
	Name        *string
	Description *string
	LastScan    *time.Time
	LastScanSet bool
}

type AssetVulnerability struct {
	ID            string    `json:"id"`
	Description   string    `json:"description"`
	Severity      Severity  `json:"severity"`
	ScanID        string    `json:"scanId"`
	ComponentID   string    `json:"componentId"`
	ComponentName string    `json:"componentName"`
	PerformedAt   time.Time `json:"performedAt"`
}

type ListAssetVulnerabilitiesResponse struct {
	Data       []AssetVulnerability `json:"data"`
	Pagination Pagination           `json:"pagination"`
}

type AssetThreat struct {
	ID            string    `json:"id"`
	Description   string    `json:"description"`
	RiskLevel     RiskLevel `json:"riskLevel"`
	Type          string    `json:"type"`
	ScanID        string    `json:"scanId"`
	ComponentID   string    `json:"componentId"`
	ComponentName string    `json:"componentName"`
	PerformedAt   time.Time `json:"performedAt"`
}

type ListAssetThreatsResponse struct {
	Data       []AssetThreat `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

type QueryValidationDetail struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
	Value string `json:"value"`
}

var ErrAssetNotFound = errors.New("asset not found")

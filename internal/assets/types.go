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

type ListAssetsQuery struct {
	NameContains string

	CreatedFrom  *time.Time
	CreatedTo    *time.Time
	LastScanFrom *time.Time
	LastScanTo   *time.Time

	Page      int
	PageSize  int
	SortBy    SortBy
	SortOrder SortOrder
}

type AssetSummary struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	CreatedAt          time.Time  `json:"createdAt"`
	LastScan           *time.Time `json:"lastScan"`
	HasVulnerabilities bool       `json:"hasVulnerabilities"`
	HasThreats         bool       `json:"hasThreats"`
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

type QueryValidationDetail struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
	Value string `json:"value"`
}

var ErrAssetNotFound = errors.New("asset not found")

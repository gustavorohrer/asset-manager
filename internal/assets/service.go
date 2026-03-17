package assets

import (
	"context"
	"math"
)

type Repository interface {
	ListAssets(ctx context.Context, query ListAssetsQuery) ([]AssetSummary, int, error)
	GetAssetDetails(ctx context.Context, assetID string) (AssetDetails, error)
	ListAssetVulnerabilities(ctx context.Context, assetID string, query ListAssetVulnerabilitiesQuery) ([]AssetVulnerability, int, error)
	ListAssetThreats(ctx context.Context, assetID string, query ListAssetThreatsQuery) ([]AssetThreat, int, error)
	UpdateAsset(ctx context.Context, assetID string, input UpdateAssetInput) (AssetUpdated, error)
}

type Lister interface {
	ListAssets(ctx context.Context, query ListAssetsQuery) (ListAssetsResponse, error)
}

type DetailsGetter interface {
	GetAssetDetails(ctx context.Context, assetID string) (AssetDetails, error)
}

type ServiceAPI interface {
	Lister
	DetailsGetter
	VulnerabilitiesLister
	ThreatsLister
	AssetUpdater
}

type VulnerabilitiesLister interface {
	ListAssetVulnerabilities(ctx context.Context, assetID string, query ListAssetVulnerabilitiesQuery) (ListAssetVulnerabilitiesResponse, error)
}

type ThreatsLister interface {
	ListAssetThreats(ctx context.Context, assetID string, query ListAssetThreatsQuery) (ListAssetThreatsResponse, error)
}

type AssetUpdater interface {
	UpdateAsset(ctx context.Context, assetID string, input UpdateAssetInput) (AssetUpdated, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListAssets(ctx context.Context, query ListAssetsQuery) (ListAssetsResponse, error) {
	data, total, err := s.repo.ListAssets(ctx, query)
	if err != nil {
		return ListAssetsResponse{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(query.PageSize)))
	}

	return ListAssetsResponse{
		Data: data,
		Pagination: Pagination{
			Page:       query.Page,
			PageSize:   query.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) GetAssetDetails(ctx context.Context, assetID string) (AssetDetails, error) {
	return s.repo.GetAssetDetails(ctx, assetID)
}

func (s *Service) ListAssetVulnerabilities(ctx context.Context, assetID string, query ListAssetVulnerabilitiesQuery) (ListAssetVulnerabilitiesResponse, error) {
	data, total, err := s.repo.ListAssetVulnerabilities(ctx, assetID, query)
	if err != nil {
		return ListAssetVulnerabilitiesResponse{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(query.PageSize)))
	}

	return ListAssetVulnerabilitiesResponse{
		Data: data,
		Pagination: Pagination{
			Page:       query.Page,
			PageSize:   query.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) ListAssetThreats(ctx context.Context, assetID string, query ListAssetThreatsQuery) (ListAssetThreatsResponse, error) {
	data, total, err := s.repo.ListAssetThreats(ctx, assetID, query)
	if err != nil {
		return ListAssetThreatsResponse{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(query.PageSize)))
	}

	return ListAssetThreatsResponse{
		Data: data,
		Pagination: Pagination{
			Page:       query.Page,
			PageSize:   query.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) UpdateAsset(ctx context.Context, assetID string, input UpdateAssetInput) (AssetUpdated, error) {
	return s.repo.UpdateAsset(ctx, assetID, input)
}

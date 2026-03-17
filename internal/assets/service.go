package assets

import (
	"context"
	"math"
)

type Repository interface {
	ListAssets(ctx context.Context, query ListAssetsQuery) ([]AssetSummary, int, error)
}

type Lister interface {
	ListAssets(ctx context.Context, query ListAssetsQuery) (ListAssetsResponse, error)
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

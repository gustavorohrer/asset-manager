package assets

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	listData  []AssetSummary
	listTotal int
	listErr   error

	detailsData AssetDetails
	detailsErr  error
}

func (f *fakeRepository) ListAssets(_ context.Context, _ ListAssetsQuery) ([]AssetSummary, int, error) {
	return f.listData, f.listTotal, f.listErr
}

func (f *fakeRepository) GetAssetDetails(_ context.Context, _ string) (AssetDetails, error) {
	return f.detailsData, f.detailsErr
}

func TestServiceListAssetsCalculatesTotalPages(t *testing.T) {
	service := NewService(&fakeRepository{
		listData:  []AssetSummary{{ID: "AST-001"}},
		listTotal: 21,
	})

	response, err := service.ListAssets(context.Background(), ListAssetsQuery{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Pagination.TotalPages != 2 {
		t.Fatalf("expected totalPages=2, got %d", response.Pagination.TotalPages)
	}
}

func TestServiceGetAssetDetailsPassThrough(t *testing.T) {
	service := NewService(&fakeRepository{
		detailsData: AssetDetails{ID: "AST-001"},
	})

	details, err := service.GetAssetDetails(context.Background(), "AST-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if details.ID != "AST-001" {
		t.Fatalf("expected id AST-001, got %s", details.ID)
	}
}

func TestServiceGetAssetDetailsError(t *testing.T) {
	service := NewService(&fakeRepository{
		detailsErr: errors.New("boom"),
	})

	_, err := service.GetAssetDetails(context.Background(), "AST-001")
	if err == nil {
		t.Fatal("expected error")
	}
}

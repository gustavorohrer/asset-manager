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

	vulnData  []AssetVulnerability
	vulnTotal int
	vulnErr   error

	threatData  []AssetThreat
	threatTotal int
	threatErr   error

	updateData  AssetUpdated
	updateErr   error
	updateID    string
	updateInput UpdateAssetInput

	deleteData AssetDeleted
	deleteErr  error
	deleteID   string
}

func (f *fakeRepository) ListAssets(_ context.Context, _ ListAssetsQuery) ([]AssetSummary, int, error) {
	return f.listData, f.listTotal, f.listErr
}

func (f *fakeRepository) GetAssetDetails(_ context.Context, _ string) (AssetDetails, error) {
	return f.detailsData, f.detailsErr
}

func (f *fakeRepository) ListAssetVulnerabilities(_ context.Context, _ string, _ ListAssetVulnerabilitiesQuery) ([]AssetVulnerability, int, error) {
	return f.vulnData, f.vulnTotal, f.vulnErr
}

func (f *fakeRepository) ListAssetThreats(_ context.Context, _ string, _ ListAssetThreatsQuery) ([]AssetThreat, int, error) {
	return f.threatData, f.threatTotal, f.threatErr
}

func (f *fakeRepository) UpdateAsset(_ context.Context, assetID string, input UpdateAssetInput) (AssetUpdated, error) {
	f.updateID = assetID
	f.updateInput = input
	return f.updateData, f.updateErr
}

func (f *fakeRepository) DeleteAsset(_ context.Context, assetID string) (AssetDeleted, error) {
	f.deleteID = assetID
	return f.deleteData, f.deleteErr
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

func TestServiceListAssetVulnerabilitiesCalculatesTotalPages(t *testing.T) {
	service := NewService(&fakeRepository{
		vulnData:  []AssetVulnerability{{ID: "VUL-001"}},
		vulnTotal: 21,
	})

	response, err := service.ListAssetVulnerabilities(context.Background(), "AST-001", ListAssetVulnerabilitiesQuery{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Pagination.TotalPages != 2 {
		t.Fatalf("expected totalPages=2, got %d", response.Pagination.TotalPages)
	}
}

func TestServiceListAssetThreatsCalculatesTotalPages(t *testing.T) {
	service := NewService(&fakeRepository{
		threatData:  []AssetThreat{{ID: "THR-001"}},
		threatTotal: 21,
	})

	response, err := service.ListAssetThreats(context.Background(), "AST-001", ListAssetThreatsQuery{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Pagination.TotalPages != 2 {
		t.Fatalf("expected totalPages=2, got %d", response.Pagination.TotalPages)
	}
}

func TestServiceUpdateAssetPassThrough(t *testing.T) {
	service := NewService(&fakeRepository{
		updateData: AssetUpdated{ID: "AST-001"},
	})

	name := "Updated"
	got, err := service.UpdateAsset(context.Background(), "AST-001", UpdateAssetInput{Name: &name})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "AST-001" {
		t.Fatalf("expected id AST-001, got %s", got.ID)
	}
}

func TestServiceDeleteAssetPassThrough(t *testing.T) {
	service := NewService(&fakeRepository{
		deleteData: AssetDeleted{
			ID:      "AST-001",
			Deleted: true,
		},
	})

	got, err := service.DeleteAsset(context.Background(), "AST-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "AST-001" || !got.Deleted {
		t.Fatalf("unexpected delete response: %+v", got)
	}
}

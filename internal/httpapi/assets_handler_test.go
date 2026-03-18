package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gustavorohrer/ecl-be-challenge/internal/assets"
)

type fakeAssetsLister struct {
	response assets.ListAssetsResponse
	err      error

	detailsResponse assets.AssetDetails
	detailsErr      error

	called        bool
	query         assets.ListAssetsQuery
	detailsCalled bool
	detailsID     string

	vulnerabilitiesResponse assets.ListAssetVulnerabilitiesResponse
	vulnerabilitiesErr      error
	vulnerabilitiesCalled   bool
	vulnerabilitiesID       string
	vulnerabilitiesQuery    assets.ListAssetVulnerabilitiesQuery

	threatsResponse assets.ListAssetThreatsResponse
	threatsErr      error
	threatsCalled   bool
	threatsID       string
	threatsQuery    assets.ListAssetThreatsQuery

	updateResponse assets.AssetUpdated
	updateErr      error
	updateCalled   bool
	updateID       string
	updateInput    assets.UpdateAssetInput

	deleteResponse assets.AssetDeleted
	deleteErr      error
	deleteCalled   bool
	deleteID       string
}

func (f *fakeAssetsLister) ListAssets(_ context.Context, query assets.ListAssetsQuery) (assets.ListAssetsResponse, error) {
	f.called = true
	f.query = query
	return f.response, f.err
}

func (f *fakeAssetsLister) GetAssetDetails(_ context.Context, assetID string) (assets.AssetDetails, error) {
	f.detailsCalled = true
	f.detailsID = assetID
	return f.detailsResponse, f.detailsErr
}

func (f *fakeAssetsLister) ListAssetVulnerabilities(_ context.Context, assetID string, query assets.ListAssetVulnerabilitiesQuery) (assets.ListAssetVulnerabilitiesResponse, error) {
	f.vulnerabilitiesCalled = true
	f.vulnerabilitiesID = assetID
	f.vulnerabilitiesQuery = query
	return f.vulnerabilitiesResponse, f.vulnerabilitiesErr
}

func (f *fakeAssetsLister) ListAssetThreats(_ context.Context, assetID string, query assets.ListAssetThreatsQuery) (assets.ListAssetThreatsResponse, error) {
	f.threatsCalled = true
	f.threatsID = assetID
	f.threatsQuery = query
	return f.threatsResponse, f.threatsErr
}

func (f *fakeAssetsLister) UpdateAsset(_ context.Context, assetID string, input assets.UpdateAssetInput) (assets.AssetUpdated, error) {
	f.updateCalled = true
	f.updateID = assetID
	f.updateInput = input
	return f.updateResponse, f.updateErr
}

func (f *fakeAssetsLister) DeleteAsset(_ context.Context, assetID string) (assets.AssetDeleted, error) {
	f.deleteCalled = true
	f.deleteID = assetID
	return f.deleteResponse, f.deleteErr
}

func TestListAssetsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	lastScan := time.Date(2024, 10, 8, 0, 0, 0, 0, time.UTC)
	lister := &fakeAssetsLister{
		response: assets.ListAssetsResponse{
			Data: []assets.AssetSummary{
				{
					ID:                 "AST-001",
					Name:               "Dell PowerEdge R740 Server",
					Description:        "Production database server",
					CreatedAt:          time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					LastScan:           &lastScan,
					HasVulnerabilities: true,
					HasThreats:         true,
				},
			},
			Pagination: assets.Pagination{
				Page:       1,
				PageSize:   20,
				Total:      1,
				TotalPages: 1,
			},
		},
	}

	handler := NewAssetsHandler(lister)
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets?page=1&pageSize=20&sortBy=name&sortOrder=asc", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if !lister.called {
		t.Fatal("expected lister to be called")
	}
	if lister.query.SortBy != assets.SortByName || lister.query.SortOrder != assets.SortOrderAsc {
		t.Fatalf("expected sortBy=name and sortOrder=asc, got sortBy=%s sortOrder=%s", lister.query.SortBy, lister.query.SortOrder)
	}

	var payload assets.ListAssetsResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(payload.Data) != 1 {
		t.Fatalf("expected one item in response data, got %d", len(payload.Data))
	}
	if payload.Pagination.Total != 1 {
		t.Fatalf("expected total=1, got %d", payload.Pagination.Total)
	}
}

func TestListAssetsInvalidQueryReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets?page=0", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	var payload errorEnvelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if payload.Error.Code != "INVALID_QUERY_PARAM" {
		t.Fatalf("expected INVALID_QUERY_PARAM, got %s", payload.Error.Code)
	}
	if len(payload.Error.Details) == 0 {
		t.Fatal("expected at least one validation detail")
	}
}

func TestGetAssetDetailsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	lastScan := time.Date(2024, 10, 8, 0, 0, 0, 0, time.UTC)
	handler := NewAssetsHandler(&fakeAssetsLister{
		detailsResponse: assets.AssetDetails{
			ID:                 "AST-001",
			Name:               "Dell PowerEdge R740 Server",
			Description:        "Production database server",
			CreatedAt:          time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			LastScan:           &lastScan,
			HasVulnerabilities: true,
			HasThreats:         true,
			Components: []assets.AssetComponent{
				{
					ID:        "CMP-001",
					Name:      "Dell UEFI BIOS",
					Version:   "2.10.2",
					Vendor:    "Dell Inc.",
					Type:      "UEFI",
					CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					LastScan:  &lastScan,
					AssetID:   "AST-001",
				},
			},
		},
	})

	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets/AST-001", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var payload assetDetailsEnvelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.Data.ID != "AST-001" {
		t.Fatalf("expected id AST-001, got %s", payload.Data.ID)
	}
	if len(payload.Data.Components) != 1 {
		t.Fatalf("expected one component, got %d", len(payload.Data.Components))
	}
}

func TestGetAssetDetailsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{detailsErr: assets.ErrAssetNotFound})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets/AST-404", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}

	var payload errorEnvelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if payload.Error.Code != "ASSET_NOT_FOUND" {
		t.Fatalf("expected ASSET_NOT_FOUND, got %s", payload.Error.Code)
	}
}

func TestGetAssetDetailsInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{detailsErr: context.DeadlineExceeded})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets/AST-001", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestListAssetVulnerabilitiesSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{
		vulnerabilitiesResponse: assets.ListAssetVulnerabilitiesResponse{
			Data: []assets.AssetVulnerability{
				{
					ID:            "VUL-001",
					Description:   "Sample vulnerability",
					Severity:      assets.SeverityCritical,
					ScanID:        "SCN-001",
					ComponentID:   "CMP-001",
					ComponentName: "Dell UEFI BIOS",
					PerformedAt:   time.Date(2024, 10, 8, 0, 0, 0, 0, time.UTC),
				},
			},
			Pagination: assets.Pagination{
				Page:       1,
				PageSize:   20,
				Total:      1,
				TotalPages: 1,
			},
		},
	})

	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets/AST-001/vulnerabilities?page=1&pageSize=20&severity=critical", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var payload assets.ListAssetVulnerabilitiesResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(payload.Data) != 1 {
		t.Fatalf("expected one vulnerability, got %d", len(payload.Data))
	}
}

func TestListAssetVulnerabilitiesInvalidQueryReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets/AST-001/vulnerabilities?severity=urgent", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestListAssetVulnerabilitiesNotFoundReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{vulnerabilitiesErr: assets.ErrAssetNotFound})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets/AST-404/vulnerabilities", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestListAssetThreatsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{
		threatsResponse: assets.ListAssetThreatsResponse{
			Data: []assets.AssetThreat{
				{
					ID:            "THR-001",
					Description:   "Sample threat",
					RiskLevel:     assets.RiskLevelHigh,
					Type:          "Firmware Implant",
					ScanID:        "SCN-001",
					ComponentID:   "CMP-001",
					ComponentName: "Dell UEFI BIOS",
					PerformedAt:   time.Date(2024, 10, 8, 0, 0, 0, 0, time.UTC),
				},
			},
			Pagination: assets.Pagination{
				Page:       1,
				PageSize:   20,
				Total:      1,
				TotalPages: 1,
			},
		},
	})

	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets/AST-001/threats?page=1&pageSize=20&riskLevel=high", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var payload assets.ListAssetThreatsResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(payload.Data) != 1 {
		t.Fatalf("expected one threat, got %d", len(payload.Data))
	}
}

func TestListAssetThreatsInvalidQueryReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets/AST-001/threats?riskLevel=critical", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestListAssetThreatsNotFoundReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{threatsErr: assets.ErrAssetNotFound})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/assets/AST-404/threats", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestUpdateAssetSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	lastScan := time.Date(2024, 10, 7, 0, 0, 0, 0, time.UTC)
	lister := &fakeAssetsLister{
		updateResponse: assets.AssetUpdated{
			ID:          "AST-001",
			Name:        "Updated asset",
			Description: "updated description",
			CreatedAt:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			LastScan:    &lastScan,
		},
	}
	handler := NewAssetsHandler(lister)
	router := gin.New()
	handler.RegisterRoutes(router)

	body := `{"name":"  Updated asset  ","description":"updated description","lastScan":"2024-10-07T00:00:00Z"}`
	request := httptest.NewRequest(http.MethodPatch, "/assets/AST-001", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if !lister.updateCalled {
		t.Fatal("expected update service to be called")
	}
	if lister.updateInput.Name == nil || *lister.updateInput.Name != "Updated asset" {
		t.Fatalf("expected trimmed name to be sent, got %#v", lister.updateInput.Name)
	}
}

func TestUpdateAssetInvalidBodyReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodPatch, "/assets/AST-001", bytes.NewBufferString(`{"unknown":"value"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	var payload errorEnvelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if payload.Error.Code != "INVALID_REQUEST_BODY" {
		t.Fatalf("expected INVALID_REQUEST_BODY, got %s", payload.Error.Code)
	}
}

func TestUpdateAssetBodyTooLargeReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{})
	router := gin.New()
	handler.RegisterRoutes(router)

	oversizedDescription := string(bytes.Repeat([]byte("a"), maxUpdateAssetRequestBodyBytes+1))
	request := httptest.NewRequest(http.MethodPatch, "/assets/AST-001", bytes.NewBufferString(`{"description":"`+oversizedDescription+`"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	var payload errorEnvelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if payload.Error.Code != "INVALID_REQUEST_BODY" {
		t.Fatalf("expected INVALID_REQUEST_BODY, got %s", payload.Error.Code)
	}
}

func TestUpdateAssetNotFoundReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{updateErr: assets.ErrAssetNotFound})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodPatch, "/assets/AST-404", bytes.NewBufferString(`{"name":"updated"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestDeleteAssetSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	lister := &fakeAssetsLister{
		deleteResponse: assets.AssetDeleted{
			ID:      "AST-001",
			Deleted: true,
		},
	}

	handler := NewAssetsHandler(lister)
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodDelete, "/assets/AST-001", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if !lister.deleteCalled {
		t.Fatal("expected delete service to be called")
	}
}

func TestDeleteAssetInvalidPathReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodDelete, "/assets/%20", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestDeleteAssetNotFoundReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAssetsHandler(&fakeAssetsLister{deleteErr: assets.ErrAssetNotFound})
	router := gin.New()
	handler.RegisterRoutes(router)

	request := httptest.NewRequest(http.MethodDelete, "/assets/AST-404", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

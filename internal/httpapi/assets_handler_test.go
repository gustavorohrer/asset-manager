package httpapi

import (
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

	called bool
	query  assets.ListAssetsQuery
}

func (f *fakeAssetsLister) ListAssets(_ context.Context, query assets.ListAssetsQuery) (assets.ListAssetsResponse, error) {
	f.called = true
	f.query = query
	return f.response, f.err
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

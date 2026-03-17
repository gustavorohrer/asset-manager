//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gustavorohrer/ecl-be-challenge/internal/assets"
	"github.com/gustavorohrer/ecl-be-challenge/internal/httpapi"
	"github.com/gustavorohrer/ecl-be-challenge/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type errorEnvelope struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type assetDetailsEnvelope struct {
	Data assets.AssetDetails `json:"data"`
}

func TestAssetsAPIIntegration(t *testing.T) {
	router, pool, cleanup := setupIntegrationRouter(t)
	defer cleanup()

	t.Run("health", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/health")
		if status != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body=%s", status, string(body))
		}
	})

	t.Run("list assets success", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/assets?page=1&pageSize=3&sortBy=createdAt&sortOrder=desc")
		if status != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body=%s", status, string(body))
		}

		var payload assets.ListAssetsResponse
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Pagination.Page != 1 || payload.Pagination.PageSize != 3 {
			t.Fatalf("unexpected pagination: %+v", payload.Pagination)
		}
		if payload.Pagination.Total <= 0 {
			t.Fatalf("expected total > 0, got %d", payload.Pagination.Total)
		}
		if len(payload.Data) == 0 {
			t.Fatal("expected non-empty data")
		}
	})

	t.Run("get asset details success", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/assets/AST-001")
		if status != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body=%s", status, string(body))
		}

		var payload assetDetailsEnvelope
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Data.ID != "AST-001" {
			t.Fatalf("expected AST-001, got %s", payload.Data.ID)
		}
		if len(payload.Data.Components) == 0 {
			t.Fatal("expected at least one component")
		}
	})

	t.Run("get asset details not found", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/assets/AST-404")
		if status != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d, body=%s", status, string(body))
		}

		var payload errorEnvelope
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if payload.Error.Code != "ASSET_NOT_FOUND" {
			t.Fatalf("expected ASSET_NOT_FOUND, got %s", payload.Error.Code)
		}
	})

	t.Run("list vulnerabilities success", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/assets/AST-001/vulnerabilities?page=1&pageSize=5&severity=critical")
		if status != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body=%s", status, string(body))
		}

		var payload assets.ListAssetVulnerabilitiesResponse
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Pagination.Page != 1 || payload.Pagination.PageSize != 5 {
			t.Fatalf("unexpected pagination: %+v", payload.Pagination)
		}
		if len(payload.Data) == 0 {
			t.Fatal("expected vulnerabilities for AST-001 with severity=critical")
		}
	})

	t.Run("list vulnerabilities invalid query", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/assets/AST-001/vulnerabilities?severity=urgent")
		if status != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d, body=%s", status, string(body))
		}

		var payload errorEnvelope
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if payload.Error.Code != "INVALID_QUERY_PARAM" {
			t.Fatalf("expected INVALID_QUERY_PARAM, got %s", payload.Error.Code)
		}
	})

	t.Run("list vulnerabilities not found", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/assets/AST-404/vulnerabilities")
		if status != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d, body=%s", status, string(body))
		}

		var payload errorEnvelope
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if payload.Error.Code != "ASSET_NOT_FOUND" {
			t.Fatalf("expected ASSET_NOT_FOUND, got %s", payload.Error.Code)
		}
	})

	t.Run("list threats success", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/assets/AST-001/threats?page=1&pageSize=5&riskLevel=high")
		if status != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body=%s", status, string(body))
		}

		var payload assets.ListAssetThreatsResponse
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Pagination.Page != 1 || payload.Pagination.PageSize != 5 {
			t.Fatalf("unexpected pagination: %+v", payload.Pagination)
		}
		if len(payload.Data) == 0 {
			t.Fatal("expected threats for AST-001 with riskLevel=high")
		}
	})

	t.Run("list threats empty result", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/assets/AST-001/threats?riskLevel=low")
		if status != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body=%s", status, string(body))
		}

		var payload assets.ListAssetThreatsResponse
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if len(payload.Data) != 0 {
			t.Fatalf("expected empty threats list, got %d items", len(payload.Data))
		}
		if payload.Pagination.Total != 0 || payload.Pagination.TotalPages != 0 {
			t.Fatalf("expected zero totals, got %+v", payload.Pagination)
		}
	})

	t.Run("list threats invalid query", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/assets/AST-001/threats?riskLevel=critical")
		if status != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d, body=%s", status, string(body))
		}

		var payload errorEnvelope
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if payload.Error.Code != "INVALID_QUERY_PARAM" {
			t.Fatalf("expected INVALID_QUERY_PARAM, got %s", payload.Error.Code)
		}
	})

	t.Run("list threats not found", func(t *testing.T) {
		status, body := performRequest(t, router, http.MethodGet, "/assets/AST-404/threats")
		if status != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d, body=%s", status, string(body))
		}

		var payload errorEnvelope
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if payload.Error.Code != "ASSET_NOT_FOUND" {
			t.Fatalf("expected ASSET_NOT_FOUND, got %s", payload.Error.Code)
		}
	})

	t.Run("update asset success", func(t *testing.T) {
		beforeStatus, beforeBody := performRequest(t, router, http.MethodGet, "/assets/AST-001")
		if beforeStatus != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body=%s", beforeStatus, string(beforeBody))
		}
		var before assetDetailsEnvelope
		if err := json.Unmarshal(beforeBody, &before); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		updatedBody := []byte(`{"name":"AST-001 Integration Updated","description":"integration test update","lastScan":"2024-10-07T00:00:00Z"}`)
		status, body := performJSONRequest(t, router, http.MethodPatch, "/assets/AST-001", updatedBody)
		if status != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body=%s", status, string(body))
		}

		var payload struct {
			Data assets.AssetUpdated `json:"data"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Data.Name != "AST-001 Integration Updated" {
			t.Fatalf("expected updated name, got %s", payload.Data.Name)
		}

		restorePayload := map[string]any{
			"name":        before.Data.Name,
			"description": before.Data.Description,
		}
		if before.Data.LastScan == nil {
			restorePayload["lastScan"] = nil
		} else {
			restorePayload["lastScan"] = before.Data.LastScan.Format(time.RFC3339)
		}
		restoreBody, err := json.Marshal(restorePayload)
		if err != nil {
			t.Fatalf("marshal restore body: %v", err)
		}
		restoreStatus, restoreResponse := performJSONRequest(t, router, http.MethodPatch, "/assets/AST-001", restoreBody)
		if restoreStatus != http.StatusOK {
			t.Fatalf("expected restore status 200, got %d, body=%s", restoreStatus, string(restoreResponse))
		}
	})

	t.Run("update asset invalid body", func(t *testing.T) {
		status, body := performJSONRequest(t, router, http.MethodPatch, "/assets/AST-001", []byte(`{"id":"AST-002"}`))
		if status != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d, body=%s", status, string(body))
		}

		var payload errorEnvelope
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if payload.Error.Code != "INVALID_REQUEST_BODY" {
			t.Fatalf("expected INVALID_REQUEST_BODY, got %s", payload.Error.Code)
		}
	})

	t.Run("update asset not found", func(t *testing.T) {
		status, body := performJSONRequest(t, router, http.MethodPatch, "/assets/AST-404", []byte(`{"name":"updated"}`))
		if status != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d, body=%s", status, string(body))
		}

		var payload errorEnvelope
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if payload.Error.Code != "ASSET_NOT_FOUND" {
			t.Fatalf("expected ASSET_NOT_FOUND, got %s", payload.Error.Code)
		}
	})

	t.Run("delete asset success and second delete not found", func(t *testing.T) {
		suffix := strconv.FormatInt(time.Now().UnixNano(), 10)
		assetID := "AST-DEL-" + suffix
		componentID := "CMP-DEL-" + suffix
		scanID := "SCN-DEL-" + suffix
		vulnID := "VUL-DEL-" + suffix
		threatID := "THR-DEL-" + suffix

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := pool.Exec(ctx, `
INSERT INTO asset (id, name, description, createdat, lastscan)
VALUES ($1, $2, $3, $4, $5)
`, assetID, "Delete Integration Asset", "asset for delete integration test", "2024-01-01", "2024-10-08"); err != nil {
			t.Fatalf("insert asset: %v", err)
		}

		if _, err := pool.Exec(ctx, `
INSERT INTO component (id, name, version, vendor, type, createdat, lastscan, assetid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`, componentID, "Delete Integration Component", "1.0.0", "Integration Vendor", "Firmware", "2024-01-01", "2024-10-08", assetID); err != nil {
			t.Fatalf("insert component: %v", err)
		}

		if _, err := pool.Exec(ctx, `
INSERT INTO scan (id, performedat, scannername, componentid)
VALUES ($1, $2, $3, $4)
`, scanID, "2024-10-08", "integration-scanner", componentID); err != nil {
			t.Fatalf("insert scan: %v", err)
		}

		if _, err := pool.Exec(ctx, `
INSERT INTO vulnerability (id, description, severity, scanid)
VALUES ($1, $2, $3, $4)
`, vulnID, "integration vuln", "CRITICAL", scanID); err != nil {
			t.Fatalf("insert vulnerability: %v", err)
		}

		if _, err := pool.Exec(ctx, `
INSERT INTO threat (id, description, risklevel, type, scanid)
VALUES ($1, $2, $3, $4, $5)
`, threatID, "integration threat", "HIGH", "Integration Threat", scanID); err != nil {
			t.Fatalf("insert threat: %v", err)
		}

		status, body := performRequest(t, router, http.MethodDelete, "/assets/"+assetID)
		if status != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body=%s", status, string(body))
		}

		var payload struct {
			Data assets.AssetDeleted `json:"data"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Data.ID != assetID || !payload.Data.Deleted {
			t.Fatalf("unexpected delete response: %+v", payload.Data)
		}

		assertCountByID(t, pool, "asset", assetID, 0)
		assertCountByID(t, pool, "component", componentID, 0)
		assertCountByID(t, pool, "scan", scanID, 0)
		assertCountByID(t, pool, "vulnerability", vulnID, 0)
		assertCountByID(t, pool, "threat", threatID, 0)

		statusSecond, bodySecond := performRequest(t, router, http.MethodDelete, "/assets/"+assetID)
		if statusSecond != http.StatusNotFound {
			t.Fatalf("expected second delete status 404, got %d, body=%s", statusSecond, string(bodySecond))
		}
	})
}

func setupIntegrationRouter(t *testing.T) (http.Handler, *pgxpool.Pool, func()) {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("integration tests require DATABASE_URL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("ping database: %v", err)
	}

	repo := postgres.NewAssetRepository(pool)
	service := assets.NewService(repo)
	handler := httpapi.NewAssetsHandler(service)
	gin.SetMode(gin.TestMode)
	router := httpapi.NewRouter(handler)

	return router, pool, pool.Close
}

func performRequest(t *testing.T, router http.Handler, method, target string) (int, []byte) {
	t.Helper()

	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec.Code, rec.Body.Bytes()
}

func performJSONRequest(t *testing.T, router http.Handler, method, target string, body []byte) (int, []byte) {
	t.Helper()

	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec.Code, rec.Body.Bytes()
}

func assertCountByID(t *testing.T, pool *pgxpool.Pool, table string, id string, expected int) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var got int
	query := "SELECT COUNT(*) FROM " + table + " WHERE id = $1"
	if err := pool.QueryRow(ctx, query, id).Scan(&got); err != nil {
		t.Fatalf("count %s by id: %v", table, err)
	}
	if got != expected {
		t.Fatalf("unexpected %s count for id=%s: want=%d got=%d", table, id, expected, got)
	}
}

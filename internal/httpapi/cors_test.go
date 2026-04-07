package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddlewarePreflightFromAllowedOrigin(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	router := newTestRouterWithCORS(t)

	req := httptest.NewRequest(http.MethodOptions, "/assets", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	req.Header.Set("Access-Control-Request-Headers", "Authorization")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("unexpected Access-Control-Allow-Origin: %q", got)
	}
	allowMethods := rec.Header().Get("Access-Control-Allow-Methods")
	if !strings.Contains(allowMethods, http.MethodGet) || !strings.Contains(allowMethods, http.MethodPatch) || !strings.Contains(allowMethods, http.MethodDelete) || !strings.Contains(allowMethods, http.MethodOptions) {
		t.Fatalf("unexpected Access-Control-Allow-Methods: %q", allowMethods)
	}
	allowHeaders := strings.ToLower(rec.Header().Get("Access-Control-Allow-Headers"))
	if !strings.Contains(allowHeaders, "accept") || !strings.Contains(allowHeaders, "content-type") || !strings.Contains(allowHeaders, "authorization") {
		t.Fatalf("unexpected Access-Control-Allow-Headers: %q", rec.Header().Get("Access-Control-Allow-Headers"))
	}
}

func TestCORSMiddlewareSimpleRequestFromAllowedOrigin(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://asset-manager-ui-pi.vercel.app/,http://localhost:3000")

	router := newTestRouterWithCORS(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://asset-manager-ui-pi.vercel.app")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://asset-manager-ui-pi.vercel.app" {
		t.Fatalf("unexpected Access-Control-Allow-Origin: %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("expected empty Access-Control-Allow-Credentials, got %q", got)
	}
}

func TestCORSMiddlewareSimpleRequestFromDisallowedOrigin(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	router := newTestRouterWithCORS(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://evil.example")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no Access-Control-Allow-Origin header, got %q", got)
	}
}

func TestNewCORSMiddlewareFromEnvFailsWhenProductionAllowlistIsEmpty(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	_, err := NewCORSMiddlewareFromEnv(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewCORSMiddlewareFromEnvFailsWhenProductionAllowlistHasWildcard(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,*")

	_, err := NewCORSMiddlewareFromEnv(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func newTestRouterWithCORS(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	handler := NewAssetsHandler(&fakeAssetsLister{})
	corsMiddleware, err := NewCORSMiddlewareFromEnv(nil)
	if err != nil {
		t.Fatalf("unexpected CORS middleware error: %v", err)
	}

	return NewRouter(handler, corsMiddleware)
}

package httpapi

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	appEnvDevelopment = "development"
	appEnvProduction  = "production"
)

func NewCORSMiddlewareFromEnv(getenv func(string) string) (gin.HandlerFunc, error) {
	if getenv == nil {
		getenv = os.Getenv
	}

	appEnv := strings.ToLower(strings.TrimSpace(getenv("APP_ENV")))
	if appEnv == "" {
		appEnv = appEnvDevelopment
	}

	allowedOrigins, err := parseCORSAllowedOrigins(getenv("CORS_ALLOWED_ORIGINS"))
	if err != nil {
		return nil, fmt.Errorf("invalid CORS_ALLOWED_ORIGINS: %w", err)
	}

	if appEnv == appEnvProduction {
		if len(allowedOrigins) == 0 {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS is required when APP_ENV=%s", appEnvProduction)
		}
		if slices.Contains(allowedOrigins, "*") {
			return nil, fmt.Errorf("wildcard origin '*' is not allowed when APP_ENV=%s", appEnvProduction)
		}
	}

	if len(allowedOrigins) == 0 {
		return nil, nil
	}

	return ginCors.New(ginCors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}), nil
}

func parseCORSAllowedOrigins(value string) ([]string, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	rawOrigins := strings.Split(value, ",")
	origins := make([]string, 0, len(rawOrigins))
	seen := make(map[string]struct{}, len(rawOrigins))

	for _, rawOrigin := range rawOrigins {
		normalizedOrigin, err := normalizeCORSOrigin(rawOrigin)
		if err != nil {
			return nil, err
		}
		if normalizedOrigin == "" {
			continue
		}
		if _, exists := seen[normalizedOrigin]; exists {
			continue
		}

		seen[normalizedOrigin] = struct{}{}
		origins = append(origins, normalizedOrigin)
	}

	return origins, nil
}

func normalizeCORSOrigin(rawOrigin string) (string, error) {
	origin := strings.TrimSpace(rawOrigin)
	if origin == "" {
		return "", nil
	}
	if origin == "*" {
		return "*", nil
	}

	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return "", fmt.Errorf("invalid origin %q", origin)
	}

	if parsedOrigin.Scheme == "" || parsedOrigin.Host == "" {
		return "", fmt.Errorf("invalid origin %q: expected scheme and host", origin)
	}
	if parsedOrigin.User != nil {
		return "", fmt.Errorf("invalid origin %q: user info is not allowed", origin)
	}
	if parsedOrigin.RawQuery != "" || parsedOrigin.Fragment != "" {
		return "", fmt.Errorf("invalid origin %q: query parameters and fragments are not allowed", origin)
	}
	if parsedOrigin.Path != "" && parsedOrigin.Path != "/" {
		return "", fmt.Errorf("invalid origin %q: path is not allowed", origin)
	}

	scheme := strings.ToLower(parsedOrigin.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("invalid origin %q: only http and https are allowed", origin)
	}

	return scheme + "://" + strings.ToLower(parsedOrigin.Host), nil
}

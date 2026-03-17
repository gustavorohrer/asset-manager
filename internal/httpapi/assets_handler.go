package httpapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gustavorohrer/ecl-be-challenge/internal/assets"
)

const requestTimeout = 5 * time.Second

type AssetsHandler struct {
	service assets.ServiceAPI
}

func NewAssetsHandler(service assets.ServiceAPI) *AssetsHandler {
	return &AssetsHandler{service: service}
}

func (h *AssetsHandler) RegisterRoutes(router gin.IRoutes) {
	router.GET("/assets", h.listAssets)
	router.GET("/assets/:id", h.getAssetDetails)
	router.GET("/assets/:id/vulnerabilities", h.listAssetVulnerabilities)
	router.GET("/assets/:id/threats", h.listAssetThreats)
	router.PATCH("/assets/:id", h.updateAsset)
	router.DELETE("/assets/:id", h.deleteAsset)
}

func (h *AssetsHandler) listAssets(c *gin.Context) {
	query, details := assets.ParseListAssetsQuery(c.Request.URL.Query())
	if len(details) > 0 {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "INVALID_QUERY_PARAM",
				Message: "one or more query parameters are invalid",
				Details: details,
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	response, err := h.service.ListAssets(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorEnvelope{
			Error: apiError{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AssetsHandler) getAssetDetails(c *gin.Context) {
	assetID := strings.TrimSpace(c.Param("id"))
	if assetID == "" {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "INVALID_PATH_PARAM",
				Message: "asset id is required",
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	response, err := h.service.GetAssetDetails(ctx, assetID)
	if err != nil {
		if errors.Is(err, assets.ErrAssetNotFound) {
			c.JSON(http.StatusNotFound, errorEnvelope{
				Error: apiError{
					Code:    "ASSET_NOT_FOUND",
					Message: "asset not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errorEnvelope{
			Error: apiError{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, assetDetailsEnvelope{Data: response})
}

func (h *AssetsHandler) listAssetVulnerabilities(c *gin.Context) {
	assetID := strings.TrimSpace(c.Param("id"))
	if assetID == "" {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "INVALID_PATH_PARAM",
				Message: "asset id is required",
			},
		})
		return
	}

	query, details := assets.ParseListAssetVulnerabilitiesQuery(c.Request.URL.Query())
	if len(details) > 0 {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "INVALID_QUERY_PARAM",
				Message: "one or more query parameters are invalid",
				Details: details,
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	response, err := h.service.ListAssetVulnerabilities(ctx, assetID, query)
	if err != nil {
		if errors.Is(err, assets.ErrAssetNotFound) {
			c.JSON(http.StatusNotFound, errorEnvelope{
				Error: apiError{
					Code:    "ASSET_NOT_FOUND",
					Message: "asset not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errorEnvelope{
			Error: apiError{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AssetsHandler) listAssetThreats(c *gin.Context) {
	assetID := strings.TrimSpace(c.Param("id"))
	if assetID == "" {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "INVALID_PATH_PARAM",
				Message: "asset id is required",
			},
		})
		return
	}

	query, details := assets.ParseListAssetThreatsQuery(c.Request.URL.Query())
	if len(details) > 0 {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "INVALID_QUERY_PARAM",
				Message: "one or more query parameters are invalid",
				Details: details,
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	response, err := h.service.ListAssetThreats(ctx, assetID, query)
	if err != nil {
		if errors.Is(err, assets.ErrAssetNotFound) {
			c.JSON(http.StatusNotFound, errorEnvelope{
				Error: apiError{
					Code:    "ASSET_NOT_FOUND",
					Message: "asset not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errorEnvelope{
			Error: apiError{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AssetsHandler) updateAsset(c *gin.Context) {
	assetID := strings.TrimSpace(c.Param("id"))
	if assetID == "" {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "INVALID_PATH_PARAM",
				Message: "asset id is required",
			},
		})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "INVALID_REQUEST_BODY",
				Message: "request body could not be read",
			},
		})
		return
	}

	input, details := assets.ParseUpdateAssetRequestBody(body)
	if len(details) > 0 {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "INVALID_REQUEST_BODY",
				Message: "request body is invalid",
				Details: details,
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	updated, err := h.service.UpdateAsset(ctx, assetID, input)
	if err != nil {
		if errors.Is(err, assets.ErrAssetNotFound) {
			c.JSON(http.StatusNotFound, errorEnvelope{
				Error: apiError{
					Code:    "ASSET_NOT_FOUND",
					Message: "asset not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errorEnvelope{
			Error: apiError{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, assetUpdatedEnvelope{Data: updated})
}

func (h *AssetsHandler) deleteAsset(c *gin.Context) {
	assetID := strings.TrimSpace(c.Param("id"))
	if assetID == "" {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "INVALID_PATH_PARAM",
				Message: "asset id is required",
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	deleted, err := h.service.DeleteAsset(ctx, assetID)
	if err != nil {
		if errors.Is(err, assets.ErrAssetNotFound) {
			c.JSON(http.StatusNotFound, errorEnvelope{
				Error: apiError{
					Code:    "ASSET_NOT_FOUND",
					Message: "asset not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errorEnvelope{
			Error: apiError{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, assetDeletedEnvelope{Data: deleted})
}

type errorEnvelope struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string                         `json:"code"`
	Message string                         `json:"message"`
	Details []assets.QueryValidationDetail `json:"details,omitempty"`
}

type assetDetailsEnvelope struct {
	Data assets.AssetDetails `json:"data"`
}

type assetUpdatedEnvelope struct {
	Data assets.AssetUpdated `json:"data"`
}

type assetDeletedEnvelope struct {
	Data assets.AssetDeleted `json:"data"`
}

package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gustavorohrer/ecl-be-challenge/internal/assets"
)

const requestTimeout = 5 * time.Second

type AssetsHandler struct {
	lister assets.Lister
}

func NewAssetsHandler(lister assets.Lister) *AssetsHandler {
	return &AssetsHandler{lister: lister}
}

func (h *AssetsHandler) RegisterRoutes(router gin.IRoutes) {
	router.GET("/assets", h.listAssets)
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

	response, err := h.lister.ListAssets(ctx, query)
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

type errorEnvelope struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string                         `json:"code"`
	Message string                         `json:"message"`
	Details []assets.QueryValidationDetail `json:"details,omitempty"`
}

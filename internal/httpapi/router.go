package httpapi

import "github.com/gin-gonic/gin"

func NewRouter(assetsHandler *AssetsHandler, middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	for _, middleware := range middlewares {
		if middleware != nil {
			router.Use(middleware)
		}
	}
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})
	assetsHandler.RegisterRoutes(router)
	return router
}

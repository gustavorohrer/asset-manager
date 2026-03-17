package httpapi

import "github.com/gin-gonic/gin"

func NewRouter(assetsHandler *AssetsHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})
	assetsHandler.RegisterRoutes(router)
	return router
}

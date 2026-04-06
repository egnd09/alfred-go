package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns server status
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

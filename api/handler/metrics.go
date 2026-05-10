package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"weather-api/metrics"
)

func Metrics(c *gin.Context) {
	c.JSON(http.StatusOK, metrics.Default.Snapshot())
}

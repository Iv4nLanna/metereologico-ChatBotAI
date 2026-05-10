package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":         "ok",
		"version":        "1.0.0",
		"uptime_seconds": int(time.Since(startTime).Seconds()),
	})
}

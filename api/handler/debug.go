package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"weather-api/static"
)

func Debug(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", static.DebugHTML)
}

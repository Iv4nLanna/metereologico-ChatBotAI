package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"weather-api/cache"
	"weather-api/service"
)

func Forecast(c *cache.Memory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		city := ctx.Query("city")
		if city == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
			return
		}

		daysStr := ctx.DefaultQuery("days", "3")
		days, err := strconv.Atoi(daysStr)
		if err != nil || days < 1 || days > 7 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "days must be between 1 and 7"})
			return
		}

		key := fmt.Sprintf("forecast:%s:%d", city, days)
		if cached, ok := c.Get(key); ok {
			ctx.JSON(http.StatusOK, cached)
			return
		}

		geo, err := service.Geocode(city)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("city not found: %s", city)})
			return
		}

		forecast, err := service.GetForecast(geo, days)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch forecast data"})
			return
		}

		c.Set(key, forecast)
		ctx.JSON(http.StatusOK, forecast)
	}
}

package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"weather-api/cache"
	"weather-api/service"
)

func Weather(c *cache.Memory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		city := ctx.Query("city")
		if city == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
			return
		}

		key := fmt.Sprintf("weather:%s", city)
		if cached, ok := c.Get(key); ok {
			ctx.JSON(http.StatusOK, cached)
			return
		}

		geo, err := service.Geocode(city)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("city not found: %s", city)})
			return
		}

		weather, err := service.GetCurrentWeather(geo)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch weather data"})
			return
		}

		c.Set(key, weather)
		ctx.JSON(http.StatusOK, weather)
	}
}

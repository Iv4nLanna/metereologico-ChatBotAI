package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"weather-api/cache"
	"weather-api/handler"
	"weather-api/middleware"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	memCache := cache.New(5 * time.Minute)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", handler.Health)
	r.GET("/weather", handler.Weather(memCache))
	r.GET("/forecast", handler.Forecast(memCache))

	slog.Info("server starting", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

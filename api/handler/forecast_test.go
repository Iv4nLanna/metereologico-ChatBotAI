package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"weather-api/cache"
)

func TestForecast_MissingCity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	c := cache.New(time.Minute)
	r.GET("/forecast", Forecast(c))

	req := httptest.NewRequest(http.MethodGet, "/forecast", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestForecast_InvalidDays(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	c := cache.New(time.Minute)
	r.GET("/forecast", Forecast(c))

	for _, d := range []string{"0", "8", "abc", "-1"} {
		req := httptest.NewRequest(http.MethodGet, "/forecast?city=SP&days="+d, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("days=%s: expected 400, got %d", d, w.Code)
		}
	}
}

func TestForecast_DefaultDays(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	c := cache.New(time.Minute)

	// Pre-populate cache to avoid real API call
	cachedData := map[string]any{
		"city": "SP",
		"days": []any{},
	}
	c.Set("forecast:SP:3", cachedData)
	r.GET("/forecast", Forecast(c))

	req := httptest.NewRequest(http.MethodGet, "/forecast?city=SP", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 with default days, got %d", w.Code)
	}
}

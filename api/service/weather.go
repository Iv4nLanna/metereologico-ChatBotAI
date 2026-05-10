package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var OpenMeteoWeatherURL = "https://api.open-meteo.com"

var weatherClient = &http.Client{Timeout: 10 * time.Second}

type CurrentWeather struct {
	City         string  `json:"city"`
	TemperatureC float64 `json:"temperature_c"`
	FeelsLikeC   float64 `json:"feels_like_c"`
	Condition    string  `json:"condition"`
	HumidityPct  int     `json:"humidity_pct"`
	WindKph      float64 `json:"wind_kph"`
	UVIndex      float64 `json:"uv_index"`
}

type ForecastDay struct {
	Date               string  `json:"date"`
	MaxTempC           float64 `json:"max_temp_c"`
	MinTempC           float64 `json:"min_temp_c"`
	Condition          string  `json:"condition"`
	PrecipitationMm    float64 `json:"precipitation_mm"`
	RainProbabilityPct int     `json:"rain_probability_pct"`
}

type Forecast struct {
	City string        `json:"city"`
	Days []ForecastDay `json:"days"`
}

func wmoCondition(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code == 1:
		return "Mainly clear"
	case code == 2:
		return "Partly cloudy"
	case code == 3:
		return "Overcast"
	case code <= 48:
		return "Foggy"
	case code <= 57:
		return "Drizzle"
	case code <= 67:
		return "Rain"
	case code <= 77:
		return "Snow"
	case code <= 82:
		return "Rain showers"
	case code <= 86:
		return "Snow showers"
	case code <= 99:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}

func GetCurrentWeather(geo *GeoResult) (*CurrentWeather, error) {
	u := fmt.Sprintf(
		"%s/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,apparent_temperature,relative_humidity_2m,wind_speed_10m,weather_code,uv_index&wind_speed_unit=kmh",
		OpenMeteoWeatherURL, geo.Latitude, geo.Longitude,
	)
	resp, err := weatherClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("weather request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var raw struct {
		Current struct {
			Temperature2m       float64 `json:"temperature_2m"`
			ApparentTemperature float64 `json:"apparent_temperature"`
			RelativeHumidity2m  int     `json:"relative_humidity_2m"`
			WindSpeed10m        float64 `json:"wind_speed_10m"`
			WeatherCode         int     `json:"weather_code"`
			UvIndex             float64 `json:"uv_index"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("weather decode failed: %w", err)
	}

	return &CurrentWeather{
		City:         geo.Name,
		TemperatureC: raw.Current.Temperature2m,
		FeelsLikeC:   raw.Current.ApparentTemperature,
		Condition:    wmoCondition(raw.Current.WeatherCode),
		HumidityPct:  raw.Current.RelativeHumidity2m,
		WindKph:      raw.Current.WindSpeed10m,
		UVIndex:      raw.Current.UvIndex,
	}, nil
}

func GetForecast(geo *GeoResult, days int) (*Forecast, error) {
	if days < 1 || days > 16 {
		return nil, fmt.Errorf("days must be between 1 and 16")
	}

	u := fmt.Sprintf(
		"%s/v1/forecast?latitude=%.4f&longitude=%.4f&daily=weather_code,temperature_2m_max,temperature_2m_min,precipitation_sum,precipitation_probability_max&forecast_days=%d&timezone=auto",
		OpenMeteoWeatherURL, geo.Latitude, geo.Longitude, days,
	)
	resp, err := weatherClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("forecast request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status %d", resp.StatusCode)
	}

	var raw struct {
		Daily struct {
			Time                        []string  `json:"time"`
			WeatherCode                 []int     `json:"weather_code"`
			Temperature2mMax            []float64 `json:"temperature_2m_max"`
			Temperature2mMin            []float64 `json:"temperature_2m_min"`
			PrecipitationSum            []float64 `json:"precipitation_sum"`
			PrecipitationProbabilityMax []int     `json:"precipitation_probability_max"`
		} `json:"daily"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("forecast decode failed: %w", err)
	}

	forecast := &Forecast{City: geo.Name}

	n := len(raw.Daily.Time)
	if n == 0 {
		return forecast, nil
	}
	for _, s := range []int{
		len(raw.Daily.WeatherCode),
		len(raw.Daily.Temperature2mMax),
		len(raw.Daily.Temperature2mMin),
		len(raw.Daily.PrecipitationSum),
		len(raw.Daily.PrecipitationProbabilityMax),
	} {
		if s < n {
			return nil, fmt.Errorf("forecast API returned inconsistent array lengths")
		}
	}

	for i := range raw.Daily.Time {
		forecast.Days = append(forecast.Days, ForecastDay{
			Date:               raw.Daily.Time[i],
			MaxTempC:           raw.Daily.Temperature2mMax[i],
			MinTempC:           raw.Daily.Temperature2mMin[i],
			Condition:          wmoCondition(raw.Daily.WeatherCode[i]),
			PrecipitationMm:    raw.Daily.PrecipitationSum[i],
			RainProbabilityPct: raw.Daily.PrecipitationProbabilityMax[i],
		})
	}
	return forecast, nil
}

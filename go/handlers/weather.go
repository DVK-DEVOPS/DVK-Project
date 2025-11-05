package handlers

import (
	"DVK-Project/client"
	"DVK-Project/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
)

type WeatherController struct {
	Client *client.APIClient
}

func NewWeatherController(c *client.APIClient) *WeatherController {
	return &WeatherController{Client: c}
}

// @Summary Serve weather page
// @Description Show the weather page.
// @Tags weather
// @Produce text/html
// @Success 200 {string} text/html "HTML of weather forecast page"
// @Failure 404 {string} string "Template not found"
// @Router /weather [get]
func (wc *WeatherController) ShowWeatherPage(w http.ResponseWriter, req *http.Request) {
	forecast, err := wc.GetForecastData(req.Context())
	if err != nil {
		// capture error with request context
		if hub := sentry.GetHubFromContext(req.Context()); hub != nil {
			hub.CaptureException(err)
		}
		http.Error(w, "Failed to get forecast", http.StatusInternalServerError)
		return
	}
	formatted := models.FormatForecastData(forecast)

	tmpl, err := template.ParseFiles("templates/weather.html")
	if err != nil {
		if hub := sentry.GetHubFromContext(req.Context()); hub != nil {
			hub.CaptureException(err)
		}
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, formatted); err != nil {
		if hub := sentry.GetHubFromContext(req.Context()); hub != nil {
			hub.CaptureException(err)
		}
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		return
	}

	if _, err := buf.WriteTo(w); err != nil {
		if hub := sentry.GetHubFromContext(req.Context()); hub != nil {
			hub.CaptureException(err)
		}
		fmt.Printf("weather.go: failed to write buffer to ResponseWriter: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (wc *WeatherController) GetForecastData(ctx context.Context) (*models.Forecast, error) {
	return wc.FetchAndParseWeatherResponse(ctx, "Copenhagen")
}

// @Summary      Get weather forecast
// @Description  Get weather forecast (temperature, conditions) for 5 days in Copenhagen
// @Produce      json
// @Tags weather
// @Success      200 {object} models.StandardResponse{data=models.Forecast}
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Router       /api/weather [get]
func (wc *WeatherController) GetWeatherForecast(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in GetWeatherForecast: %v", r)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	log.Println("GET /api/weather called")

	w.Header().Set("Content-Type", "application/json")
	forecast, err := wc.GetForecastData(req.Context())
	if err != nil {
		log.Printf("fetch error: %v", err)
		if hub := sentry.GetHubFromContext(req.Context()); hub != nil {
			hub.CaptureException(err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encErr != nil {
			if hub := sentry.GetHubFromContext(req.Context()); hub != nil {
				hub.CaptureException(encErr)
			}
			fmt.Printf("failed to write JSON error response: %v\n", encErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("GetWeatherForecast: Forecast fetched: %+v\n", forecast)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(forecast); err != nil {
		if hub := sentry.GetHubFromContext(req.Context()); hub != nil {
			hub.CaptureException(err)
		}
		fmt.Printf("failed to write JSON response: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (wc *WeatherController) FetchAndParseWeatherResponse(ctx context.Context, city string) (*models.Forecast, error) {
	log.Printf("Fetching weather for city: %s\n", city)

	// Trace external API call
	span := sentry.StartSpan(ctx, "http.client", sentry.WithDescription("FetchForecast"))
	data, err := wc.Client.FetchForecast(city)
	span.Finish()
	if err != nil {
		log.Printf("FetchAndParseWeatherResponse FetchForecast error: %v\n", err)
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.CaptureException(err)
		}
		return nil, fmt.Errorf("failed to fetch forecast: %w", err)
	}
	log.Println("FetchAndParseWeatherResponse: FetchForecast succeeded, data length:", len(data))

	// Trace parsing time
	parseSpan := sentry.StartSpan(ctx, "parse.response", sentry.WithDescription("ParseApiResponse"))
	forecast, err := models.ParseApiResponse(data)
	parseSpan.Finish()
	if err != nil {
		log.Printf("FetchAndParseWeatherResponse: ParseApiResponse error: %v\n", err)
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.CaptureException(err)
		}
		return nil, fmt.Errorf("failed to parse forecast: %w", err)
	}

	if forecast != nil && len(forecast.List) > 0 {
		log.Printf("FetchAndParseWeatherResponse: First forecast temp: %v\n", forecast.List[0].Main.Temp)
	} else {
		log.Println("FetchAndParseWeatherResponse: Forecast list is empty or nil")
	}

	log.Println("FetchAndParseWeatherResponse returning successfully")
	return forecast, nil
}

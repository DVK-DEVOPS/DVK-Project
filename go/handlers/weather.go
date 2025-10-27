package handlers

import (
	"DVK-Project/client"
	"DVK-Project/models"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
	tmpl, _ := template.ParseFiles("templates/weather.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
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
	forecast, err := wc.FetchAndParseWeatherResponse("Copenhagen")
	if err != nil {
		log.Printf("fetch error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("GetWeatherForecast: Forecast fetched: %+v\n", forecast)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(forecast)
}

func (wc *WeatherController) FetchAndParseWeatherResponse(city string) (*models.Forecast, error) {
	log.Printf("Fetching weather for city: %s\n", city)

	data, err := wc.Client.FetchForecast(city)
	if err != nil {
		log.Printf("FetchAndParseWeatherResponse FetchForecast error: %v\n", err)
		return nil, fmt.Errorf("failed to fetch forecast: %w", err)
	}
	log.Println("FetchAndParseWeatherResponse: FetchForecast succeeded, data length:", len(data))

	forecast, err := models.ParseApiResponse(data)
	if err != nil {
		log.Printf("FetchAndParseWeatherResponse: ParseApiResponse error: %v\n", err)
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

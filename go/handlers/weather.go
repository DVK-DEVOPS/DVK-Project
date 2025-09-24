package handlers

import (
	"DVK-Project/client"
	"DVK-Project/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type WeatherController struct {
	Client *client.APIClient
}

func NewWeatherController(c *client.APIClient) *WeatherController {
	return &WeatherController{Client: c}
}

func (wc *WeatherController) ShowWeatherPage(w http.ResponseWriter, req *http.Request) {
	tmpl, _ := template.ParseFiles("templates/weather.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
}

func (wc *WeatherController) GetWeatherForecast(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	forecast, err := wc.FetchAndParseWeatherResponse("Copenhagen")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.StandardResponse{
			Data: map[string]interface{}{
				"error": "failed to fetch weather forecast",
			},
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.StandardResponse{
		Data: map[string]interface{}{
			"data": forecast,
		},
	})
}

func (wc *WeatherController) FetchAndParseWeatherResponse(city string) (*models.Forecast, error) {
	data, err := wc.Client.FetchForecast(city)
	fmt.Println(data) //debug
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast: %w", err)
	}

	forecast, err := models.ParseApiResponse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse forecast: %w", err)
	}

	fmt.Println(forecast.List[0].Main.Temp) //for debug
	return forecast, nil
}

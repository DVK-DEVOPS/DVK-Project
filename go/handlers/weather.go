package handlers

import (
	"DVK-Project/models"
	"encoding/json"
	"html/template"
	"net/http"
)

type WeatherController struct {
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.StandardResponse{
		Data: nil,
	})
}

//func FetchForecast(url string) ([]byte, error) {
//	client := models.NewAPIClient()
//}

package handlers

import (
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

}

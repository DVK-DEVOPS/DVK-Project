package models

import (
	"encoding/json"
	"strings"
	"time"
)

type Forecast struct {
	List []struct {
		DtText string `json:"dt_txt"`
		Main   struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
	} `json:"list"`
}

func ParseApiResponse(data []byte) (*Forecast, error) {
	var forecast Forecast
	err := json.Unmarshal(data, &forecast)
	if err != nil {
		return nil, err
	}
	return &forecast, nil
}

type ForecastDto struct {
	Date        string
	Temperature float64
	Description string
}

func FormatForecastData(forecast *Forecast) []ForecastDto {
	formatted := []ForecastDto{}

	for _, item := range forecast.List {
		if strings.HasSuffix(item.DtText, "12:00:00") {
			t, err := time.Parse("2006-01-02 15:04:05", item.DtText)
			if err != nil {
				continue
			}

			formatted = append(formatted, ForecastDto{
				Date:        t.Format("02-01-2006"),
				Temperature: item.Main.Temp - 273.15,
				Description: item.Weather[0].Description,
			})
		}
	}
	return formatted
}

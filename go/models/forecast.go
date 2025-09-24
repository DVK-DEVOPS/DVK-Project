package models

import "encoding/json"

type Forecast struct {
	List []struct {
		DtText string `json:"dt_text"`
		Main   struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
	} `json:"list"`
}

func ParseForecast(data []byte) (*Forecast, error) {
	var forecast Forecast
	err := json.Unmarshal(data, &forecast)
	if err != nil {
		return nil, err
	}
	return &forecast, nil
}

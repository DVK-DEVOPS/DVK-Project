package models

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

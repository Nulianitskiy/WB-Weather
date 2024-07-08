package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wb-weather/internal/model"
)

type WeatherExternal interface {
	FetchForecast(lat, lon string) (*model.WeatherData, error)
}

type weatherExternal struct {
	apiKey string
}

func NewWeatherExternal(key string) WeatherExternal {
	return &weatherExternal{apiKey: key}
}

func (w *weatherExternal) FetchForecast(lat, lon string) (*model.WeatherData, error) {
	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%s&lon=%s&units=metric&appid=%s", lat, lon, w.apiKey)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weather model.WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}

	return &weather, nil
}

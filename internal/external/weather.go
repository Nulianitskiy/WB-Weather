package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wb-weather/internal/model"
)

type WeatherExternal interface {
	FetchForecast(lat, lon string) ([]model.Weather, error)
}

type weatherExternal struct {
	apiKey string
}

func NewWeatherExternal(key string) WeatherExternal {
	return &weatherExternal{apiKey: key}
}

func (w *weatherExternal) FetchForecast(lat, lon string) ([]model.Weather, error) {
	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%s&lon=%s&units=metric&appid=%s", lat, lon, w.apiKey)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp model.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var weathers []model.Weather
	for _, item := range apiResp.List {
		data, _ := json.Marshal(item)
		weather := model.Weather{
			Temp:        item.Main.Temp,
			Date:        item.DtTxt,
			WeatherData: data,
		}
		weathers = append(weathers, weather)
	}

	return weathers, nil
}

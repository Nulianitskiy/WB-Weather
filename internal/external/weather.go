package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wb-weather/internal/model"
	"wb-weather/pkg/logger"
)

type WeatherExternal interface {
	FetchForecast(lat, lon string) ([]model.Weather, error)
}

type weatherExternal struct {
	apiKey string
	log    logger.Logger
}

func NewWeatherExternal(key string, log logger.Logger) WeatherExternal {
	return &weatherExternal{apiKey: key, log: log}
}

func (w *weatherExternal) FetchForecast(lat, lon string) ([]model.Weather, error) {
	w.log.Info("Запрос прогноза погоды", "latitude", lat, "longitude", lon)
	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%s&lon=%s&units=metric&appid=%s", lat, lon, w.apiKey)

	w.log.Info("URL для запроса прогноза погоды", "url", apiUrl)
	resp, err := http.Get(apiUrl)
	if err != nil {
		w.log.Error("Ошибка при выполнении запроса к внешнему API", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp model.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		w.log.Error("Ошибка при декодировании ответа от внешнего API", "error", err)
		return nil, err
	}

	var weathers []model.Weather
	for _, item := range apiResp.List {
		data, err := json.Marshal(item)
		if err != nil {
			w.log.Error("Ошибка при маршалинге данных погоды", "error", err)
			continue
		}
		weather := model.Weather{
			Temp:        item.Main.Temp,
			Date:        item.DtTxt,
			WeatherData: data,
		}
		weathers = append(weathers, weather)
	}

	w.log.Info("Прогноз погоды успешно получен", "count", len(weathers))
	return weathers, nil
}

package repository

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"wb-weather/internal/model"
	"wb-weather/pkg/logger"
)

type WeatherRepo interface {
	UpdateCityWeatherBrute(c model.City, weathers model.WeatherData) error
	GetWeatherById(weatherId int) (model.JSONBWeather, error)
	GetForecastByCity(cityId int) ([]model.ResponseWeather, error)
}

type weatherRepo struct {
	db *sqlx.DB
}

func NewWeatherRepo(db *sqlx.DB) WeatherRepo {
	return &weatherRepo{db: db}
}

// UpdateCityWeatherBrute простое решение обновления данных
func (w *weatherRepo) UpdateCityWeatherBrute(c model.City, weathers model.WeatherData) error {
	// Удаляем прошлые данные, добавляем новые
	// TODO Придумать решение поизящней

	var oldWeather []model.ResponseWeather
	err := w.db.Select(&oldWeather, "SELECT id FROM weather WHERE city_id = $1", c.Id)
	if err != nil {
		return err
	}

	query := `DELETE FROM weather WHERE id = $1`
	for _, weather := range oldWeather {
		_, err := w.db.Exec(query, weather.Id)
		if err != nil {
			logger.Error("Ошибка при удалении информации о сотруднике", zap.Error(err), zap.Int("id", weather.Id))
			return err
		}
	}

	query = `INSERT INTO weather (city_id, date, temperature, weather_data) VALUES ($1, $2, $3, $4)`
	for _, weather := range weathers.List {
		rawJSON, err := json.Marshal(weather)
		if err != nil {
			logger.Error("Ошибка добавления данных в базу JSON", zap.Error(err))
			return err
		}
		_, err = w.db.Exec(query, c.Id, weather.DtTxt, weather.Main.Temp, string(rawJSON))
		if err != nil {
			logger.Error("Ошибка при добавлении погоды", zap.Error(err))
			return err
		}
	}
	return nil
}

func (w *weatherRepo) GetWeatherById(weatherId int) (model.JSONBWeather, error) {
	var weather model.JSONBWeather
	var jsonData []byte
	err := w.db.Get(&jsonData, "SELECT weather_data FROM weather WHERE id = $1", weatherId)
	if err != nil {
		return weather, err
	}

	err = json.Unmarshal(jsonData, &weather.Data)
	if err != nil {
		return weather, err
	}

	return weather, nil
}

func (w *weatherRepo) GetForecastByCity(cityId int) ([]model.ResponseWeather, error) {
	var forecast []model.ResponseWeather
	err := w.db.Select(&forecast, "SELECT id, city_id, date, temperature FROM weather WHERE city_id = $1", cityId)
	if err != nil {
		return forecast, err
	}
	return forecast, nil
}

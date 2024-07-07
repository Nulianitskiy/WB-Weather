package repository

import (
	"encoding/json"
	"go.uber.org/zap"
	"wb-weather/internal/model"
	"wb-weather/pkg/logger"
)

func (d *Database) GetAllCityWeather(c model.City) ([]model.ResponseWeather, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var weathers []model.ResponseWeather
	err := d.db.Select(&weathers, "SELECT id, city_id, date, temperature FROM weather WHERE city_id = $1", c.Id)
	if err != nil {
		return nil, err
	}
	return weathers, nil
}

// WeatherUpdateBrute простое решение обновления данных
func (d *Database) WeatherUpdateBrute(c model.City, weathers model.WeatherData) error {
	// Удаляем прошлые данные, добавляем новые
	// TODO Придумать решение поизящней
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var oldWeather []model.ResponseWeather
	err := d.db.Select(&oldWeather, "SELECT id FROM weather WHERE city_id = $1", c.Id)
	if err != nil {
		return err
	}

	query := `DELETE FROM weather WHERE id = $1`
	for _, weather := range oldWeather {
		_, err := d.db.Exec(query, weather.Id)
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
		_, err = d.db.Exec(query, c.Id, weather.DtTxt, weather.Main.Temp, string(rawJSON))
		if err != nil {
			logger.Error("Ошибка при добавлении погоды", zap.Error(err))
			return err
		}
	}
	return nil
}

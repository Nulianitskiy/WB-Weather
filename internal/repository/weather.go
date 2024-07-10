package repository

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
	"wb-weather/internal/model"
	"wb-weather/pkg/logger"
)

type WeatherRepo interface {
	UpdateCityWeather(c model.City, weathers []model.Weather) error
	GetWeatherByCityAndDate(city, date string) (model.ResponseFullWeather, error)
	GetForecastByCity(city string) (model.ResponseShortWeatherByCity, error)
}

type weatherRepo struct {
	db *sqlx.DB
}

func NewWeatherRepo(db *sqlx.DB) WeatherRepo {
	return &weatherRepo{db: db}
}

func (w *weatherRepo) UpdateCityWeather(c model.City, weathers []model.Weather) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	deleteQuery := `DELETE FROM weather WHERE city_id = $1 AND date < $2`
	result, err := w.db.Exec(deleteQuery, c.Id, currentTime)
	if err != nil {
		logger.Error("Ошибка при удалении старых записей", zap.Error(err))
		return err
	}
	rowsDeleted, _ := result.RowsAffected()
	logger.Info(fmt.Sprintf("Удалено старых записей: %d", rowsDeleted))

	upsertQuery := `
		INSERT INTO weather (city_id, date, temperature, weather_data)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (city_id, date) DO UPDATE SET
		temperature = EXCLUDED.temperature,
		weather_data = EXCLUDED.weather_data;
	`

	for _, weather := range weathers {
		_, err = w.db.Exec(upsertQuery, c.Id, weather.Date, weather.Temp, weather.WeatherData)
		if err != nil {
			logger.Error("Ошибка при добавлении погоды", zap.Error(err))
			return err
		}
	}

	logger.Info("Погода успешно обновлена")
	return nil
}

func (w *weatherRepo) GetWeatherByCityAndDate(city, date string) (model.ResponseFullWeather, error) {
	var weather model.ResponseFullWeather

	query := `
		SELECT 
			w.weather_data
		FROM 
			city c
		JOIN 
			weather w ON c.id = w.city_id
		WHERE 
			c.name = $1
			AND w.date::text = $2;
	`

	var weatherData []byte
	err := w.db.Get(&weatherData, query, city, date)
	if err != nil {
		return weather, err
	}

	err = json.Unmarshal(weatherData, &weather.WeatherData)
	if err != nil {
		return weather, err
	}

	return weather, nil
}

func (w *weatherRepo) GetForecastByCity(city string) (model.ResponseShortWeatherByCity, error) {
	var weather model.ResponseShortWeatherByCity

	query := `
		SELECT 
			c.name,
			c.country,
			ROUND(AVG(w.temperature),2) AS temp
		FROM 
			city c
		JOIN 
			weather w ON c.id = w.city_id
		WHERE 
			c.name = $1
		GROUP BY 
			c.id;
	`

	err := w.db.Get(&weather, query, city)
	if err != nil {
		return weather, err
	}

	var dates []string
	query = `
		SELECT DISTINCT date::text 
        FROM weather 
        WHERE city_id = (SELECT id FROM city WHERE name = $1)
		ORDER BY date
	`
	err = w.db.Select(&dates, query, city)
	if err != nil {
		return weather, err
	}
	fmt.Println(dates)
	weather.Date = dates

	return weather, nil
}

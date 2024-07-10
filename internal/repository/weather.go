package repository

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
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
	db  *sqlx.DB
	log logger.Logger
}

func NewWeatherRepo(db *sqlx.DB, log logger.Logger) WeatherRepo {
	return &weatherRepo{db: db, log: log}
}

func (w *weatherRepo) UpdateCityWeather(c model.City, weathers []model.Weather) error {
	w.log.Info("Начало обновления погоды для города", "city", c.Name)

	tx, err := w.db.Begin()
	if err != nil {
		w.log.Error("Ошибка начала транзакции", "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				w.log.Error("Ошибка отката транзакции", "error", rollbackErr)
			}
			w.log.Error("Транзакция не выполнена", "error", err)
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				w.log.Error("Ошибка фиксации транзакции", "error", commitErr)
				err = commitErr
			}
		}
	}()

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	deleteQuery := `DELETE FROM weather WHERE city_id = $1 AND date < $2`
	result, err := tx.Exec(deleteQuery, c.Id, currentTime)
	if err != nil {
		w.log.Error("Ошибка при удалении старых записей", "error", err)
		return err
	}
	rowsDeleted, _ := result.RowsAffected()
	w.log.Info("Удалено старых записей", "city", c.Name, "count", rowsDeleted)

	upsertQuery := `
		INSERT INTO weather (city_id, date, temperature, weather_data)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (city_id, date) DO UPDATE SET
		temperature = EXCLUDED.temperature,
		weather_data = EXCLUDED.weather_data;
	`

	for _, weather := range weathers {
		_, err = tx.Exec(upsertQuery, c.Id, weather.Date, weather.Temp, weather.WeatherData)
		if err != nil {
			w.log.Error("Ошибка при добавлении погоды", "error", err)
			return err
		}
	}

	sequenceResetQuery := `
		SELECT setval(pg_get_serial_sequence('weather', 'id'), coalesce(max(id), 1))
		FROM weather;
	`
	_, err = tx.Exec(sequenceResetQuery)
	if err != nil {
		w.log.Error("Ошибка сброса последовательности", "error", err)
		return err
	}

	w.log.Info("Погода успешно обновлена для города", "city", c.Name)
	return nil
}

func (w *weatherRepo) GetWeatherByCityAndDate(city, date string) (model.ResponseFullWeather, error) {
	w.log.Info("Запрос на получение погоды по городу и дате", "city", city, "date", date)
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
		w.log.Error("Ошибка при получении погоды по городу и дате", "error", err)
		return weather, err
	}

	err = json.Unmarshal(weatherData, &weather.WeatherData)
	if err != nil {
		w.log.Error("Ошибка при разборе данных погоды", "error", err)
		return weather, err
	}

	w.log.Info("Погода успешно получена", "city", city, "date", date)
	return weather, nil
}

func (w *weatherRepo) GetForecastByCity(city string) (model.ResponseShortWeatherByCity, error) {
	w.log.Info("Запрос на получение прогноза погоды по городу", "city", city)
	var weather model.ResponseShortWeatherByCity

	query := `
		SELECT 
			c.name,
			c.country,
			ROUND(AVG(w.temperature), 2) AS temp
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
		w.log.Error("Ошибка при получении прогноза погоды по городу", "error", err)
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
		w.log.Error("Ошибка при получении дат прогноза погоды", "error", err)
		return weather, err
	}
	weather.Date = dates

	w.log.Info("Прогноз погоды успешно получен", "city", city)
	return weather, nil
}

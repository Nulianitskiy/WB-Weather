package repository

import (
	"github.com/jmoiron/sqlx"
	"wb-weather/internal/model"
	"wb-weather/pkg/logger"
)

type CityRepo interface {
	AddCity(c model.City) (model.City, error)
	GetAllCitiesWithForecast() (model.CityList, error)
	GetAllCities() ([]model.City, error)
}

type cityRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewCityRepo(db *sqlx.DB, log logger.Logger) CityRepo {
	return &cityRepo{db: db, log: log}
}

// AddCity добавление города
func (r *cityRepo) AddCity(c model.City) (model.City, error) {
	r.log.Info("Начало добавления города", "name", c.Name, "country", c.Country)
	query := `INSERT INTO city (name, country, latitude, longitude) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.db.QueryRow(query, c.Name, c.Country, c.Latitude, c.Longitude).Scan(&id)
	if err != nil {
		r.log.Error("Ошибка при добавлении города", "error", err)
		return c, err
	}
	c.Id = id
	r.log.Info("Город успешно добавлен", "id", id)
	return c, nil
}

// GetAllCitiesWithForecast Получить все города
func (r *cityRepo) GetAllCitiesWithForecast() (model.CityList, error) {
	r.log.Info("Начало получения всех городов с предсказанием погоды")
	var cityNames []string
	query := `
		SELECT DISTINCT p.name
		FROM city p
		JOIN weather w ON p.id = w.city_id
		ORDER BY p.name
	`
	err := r.db.Select(&cityNames, query)
	if err != nil {
		r.log.Error("Ошибка при получении всех городов с предсказанием погоды", "error", err)
		return model.CityList{}, err
	}
	r.log.Info("Успешное получение всех городов с предсказанием погоды", "count", len(cityNames))
	return model.CityList{Cities: cityNames}, nil
}

func (r *cityRepo) GetAllCities() ([]model.City, error) {
	r.log.Info("Начало получения всех городов")
	var cities []model.City

	query := `
		SELECT * FROM city 
	`
	err := r.db.Select(&cities, query)
	if err != nil {
		r.log.Error("Ошибка при получении всех городов", "error", err)
		return nil, err
	}
	r.log.Info("Успешное получение всех городов", "count", len(cities))
	return cities, nil
}

package repository

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"wb-weather/internal/model"
	"wb-weather/pkg/logger"
)

type CityRepo interface {
	AddCity(c model.City) (model.City, error)
	GetAllCitiesWithForecast() (model.CityList, error)
	GetAllCities() ([]model.City, error)
}

type cityRepo struct {
	db *sqlx.DB
}

func NewCityRepo(db *sqlx.DB) CityRepo {
	return &cityRepo{db: db}
}

// AddCity добавление города
func (r *cityRepo) AddCity(c model.City) (model.City, error) {
	query := `INSERT INTO city (name, country, latitude, longitude) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.db.QueryRow(query, c.Name, c.Country, c.Latitude, c.Longitude).Scan(&id)
	if err != nil {
		logger.Error("Ошибка при добавлении города", zap.Error(err))
		return c, err
	}
	c.Id = id
	logger.Info("Город успешно добавлен", zap.Int("id", id))
	return c, nil
}

// GetAllCitiesWithForecast Получить все города
func (r *cityRepo) GetAllCitiesWithForecast() (model.CityList, error) {
	var cityNames []string
	// Уникальные города, которые имеют предсказание погоды в отсортированном порядке
	query := `
		SELECT DISTINCT p.name
		FROM city p
		JOIN weather w ON p.id = w.city_id
		ORDER BY p.name
	`
	err := r.db.Select(&cityNames, query)
	if err != nil {
		return model.CityList{}, err
	}
	return model.CityList{Cities: cityNames}, nil
}

func (r *cityRepo) GetAllCities() ([]model.City, error) {
	var cities []model.City

	query := `
		SELECT * FROM city 
	`
	err := r.db.Select(&cities, query)
	if err != nil {
		return nil, err
	}
	return cities, nil
}

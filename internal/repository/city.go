package repository

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"wb-weather/internal/model"
	"wb-weather/pkg/logger"
)

type CityRepo interface {
	AddCity(c model.City) (model.City, error)
	GetAllCity() ([]model.City, error)
}

type cityRepo struct {
	db *sqlx.DB
}

func NewCityRepo(db *sqlx.DB) CityRepo {
	return &cityRepo{db: db}
}

// AddCity добавление города
func (r *cityRepo) AddCity(c model.City) (model.City, error) {
	query := `INSERT INTO place (city, country, latitude, longitude) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.db.QueryRow(query, c.City, c.Country, c.Latitude, c.Longitude).Scan(&id)
	if err != nil {
		logger.Error("Ошибка при добавлении города", zap.Error(err))
		return c, err
	}
	c.Id = id
	logger.Info("Город успешно добавлен", zap.Int("id", id))
	return c, nil
}

// GetAllCity Получить все города
func (r *cityRepo) GetAllCity() ([]model.City, error) {
	var cities []model.City
	err := r.db.Select(&cities, "SELECT * FROM place ORDER BY city")
	if err != nil {
		return nil, err
	}
	return cities, nil
}

package repository

import (
	"go.uber.org/zap"
	"wb-weather/internal/model"
	"wb-weather/pkg/logger"
)

// AddCity добавление города
func (d *Database) AddCity(c model.City) (model.City, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	query := `INSERT INTO place (city, country, latitude, longitude) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := d.db.QueryRow(query, c.City, c.Country, c.Latitude, c.Longitude).Scan(&id)
	if err != nil {
		logger.Error("Ошибка при добавлении города", zap.Error(err))
		return c, err
	}
	c.Id = id
	logger.Info("Город успешно добавлен", zap.Int("id", id))
	return c, nil
}

// GetAllCity Получить все города
func (d *Database) GetAllCity() ([]model.City, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var cities []model.City
	err := d.db.Select(&cities, "SELECT * FROM place")
	if err != nil {
		return nil, err
	}
	return cities, nil
}

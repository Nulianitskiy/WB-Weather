package service

import (
	"wb-weather/internal/model"
	"wb-weather/internal/repository"
	"wb-weather/pkg/logger"
)

type CityService interface {
	AddCity(c model.City) (model.City, error)
	GetCities() (model.CityList, error)
}

type cityService struct {
	repo repository.CityRepo
	log  logger.Logger
}

func NewCityService(repo repository.CityRepo, log logger.Logger) CityService {
	return &cityService{repo: repo, log: log}
}

func (s *cityService) AddCity(c model.City) (model.City, error) {
	s.log.Info("Начало добавления города в сервисе", "name", c.Name, "country", c.Country)
	city, err := s.repo.AddCity(c)
	if err != nil {
		s.log.Error("Ошибка при добавлении города в сервисе", "error", err)
		return model.City{}, err
	}
	s.log.Info("Город успешно добавлен в сервисе", "id", city.Id)
	return city, nil
}

func (s *cityService) GetCities() (model.CityList, error) {
	s.log.Info("Начало получения всех городов с предсказанием погоды в сервисе")
	cities, err := s.repo.GetAllCitiesWithForecast()
	if err != nil {
		s.log.Error("Ошибка при получении всех городов с предсказанием погоды в сервисе", "error", err)
		return model.CityList{}, err
	}
	s.log.Info("Успешное получение всех городов с предсказанием погоды в сервисе", "count", len(cities.Cities))
	return cities, nil
}

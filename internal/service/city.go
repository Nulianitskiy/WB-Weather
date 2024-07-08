package service

import (
	"wb-weather/internal/model"
	"wb-weather/internal/repository"
)

type CityService interface {
	AddCity(c model.City) (model.City, error)
	GetCity() ([]model.City, error)
}

type cityService struct {
	repo repository.CityRepo
}

func NewCityService(repo repository.CityRepo) CityService {
	return &cityService{repo: repo}
}

func (s *cityService) AddCity(c model.City) (model.City, error) {
	return s.repo.AddCity(c)
}

func (s *cityService) GetCity() ([]model.City, error) {
	return s.repo.GetAllCity()
}

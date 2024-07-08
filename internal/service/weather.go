package service

import (
	"go.uber.org/zap"
	"wb-weather/internal/external"
	"wb-weather/internal/model"
	"wb-weather/internal/repository"
	"wb-weather/pkg/logger"
)

type WeatherService interface {
	UpdateCityWeather() error
	GetWeatherById(weatherId int) (model.JSONBWeather, error)
	GetForecastByCity(cityId int) ([]model.ResponseWeather, error)
}

type weatherService struct {
	wRepo repository.WeatherRepo
	cRepo repository.CityRepo
	ext   external.WeatherExternal
}

func NewWeatherService(wRepo repository.WeatherRepo, cRepo repository.CityRepo, ext external.WeatherExternal) WeatherService {
	return &weatherService{wRepo: wRepo, cRepo: cRepo, ext: ext}
}

func (s *weatherService) UpdateCityWeather() error {
	cities, err := s.cRepo.GetAllCity()
	if err != nil {
		logger.Error("Ошибка получения городов из бд", zap.Error(err))
		return err
	}

	for _, city := range cities {
		weatherData, err := s.ext.FetchForecast(city.Latitude, city.Longitude)
		if err != nil {
			logger.Error("Ошибка получения экземпляра базы данных", zap.Error(err))
			return err
		}
		err = s.wRepo.UpdateCityWeatherBrute(city, *weatherData)
		if err != nil {
			logger.Error("Ошибка получения экземпляра базы данных", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *weatherService) GetWeatherById(weatherId int) (model.JSONBWeather, error) {
	return s.wRepo.GetWeatherById(weatherId)
}

func (s *weatherService) GetForecastByCity(cityId int) ([]model.ResponseWeather, error) {
	return s.wRepo.GetForecastByCity(cityId)
}

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
	GetWeather(city, date string) (model.ResponseFullWeather, error)
	GetForecast(city string) (model.ResponseShortWeatherByCity, error)
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
	cities, err := s.cRepo.GetAllCities()
	if err != nil {
		logger.Error("Ошибка получения городов из бд", zap.Error(err))
		return err
	}

	// TODO распараллелить
	for _, city := range cities {
		weatherData, err := s.ext.FetchForecast(city.Latitude, city.Longitude)
		if err != nil {
			logger.Error("Ошибка получения экземпляра базы данных", zap.Error(err))
			return err
		}
		err = s.wRepo.UpdateCityWeather(city, weatherData)
		if err != nil {
			logger.Error("Ошибка получения экземпляра базы данных", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *weatherService) GetWeather(city, date string) (model.ResponseFullWeather, error) {
	return s.wRepo.GetWeatherByCityAndDate(city, date)
}

func (s *weatherService) GetForecast(city string) (model.ResponseShortWeatherByCity, error) {
	return s.wRepo.GetForecastByCity(city)
}

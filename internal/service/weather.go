package service

import (
	"fmt"
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

	errCh := make(chan error)
	doneCh := make(chan struct{})

	for _, city := range cities {
		go updateSingleCityWeather(s, city, errCh)
	}

	go func() {
		defer close(doneCh)
		for range cities {
			if err := <-errCh; err != nil {
				logger.Error(err.Error())
			}
		}
	}()
	<-doneCh
	return nil
}

func updateSingleCityWeather(s *weatherService, city model.City, errCh chan<- error) {
	weatherData, err := s.ext.FetchForecast(city.Latitude, city.Longitude)
	if err != nil {
		errCh <- fmt.Errorf("ошибка получения прогноза для города %s: %w", city.Name, err)
		return
	}
	err = s.wRepo.UpdateCityWeather(city, weatherData)
	if err != nil {
		errCh <- fmt.Errorf("ошибка обновления погоды для города %s: %w", city.Name, err)
		return
	}
	errCh <- nil
}

func (s *weatherService) GetWeather(city, date string) (model.ResponseFullWeather, error) {
	return s.wRepo.GetWeatherByCityAndDate(city, date)
}

func (s *weatherService) GetForecast(city string) (model.ResponseShortWeatherByCity, error) {
	return s.wRepo.GetForecastByCity(city)
}

package service

import (
	"fmt"
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
	log   logger.Logger
}

func NewWeatherService(wRepo repository.WeatherRepo, cRepo repository.CityRepo, ext external.WeatherExternal, log logger.Logger) WeatherService {
	return &weatherService{wRepo: wRepo, cRepo: cRepo, ext: ext, log: log}
}

func (s *weatherService) UpdateCityWeather() error {
	s.log.Info("Начало обновления погоды для всех городов")
	cities, err := s.cRepo.GetAllCities()
	if err != nil {
		s.log.Error("Ошибка получения городов из БД", "error", err)
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
				s.log.Error(err.Error())
			}
		}
	}()
	<-doneCh
	s.log.Info("Завершение обновления погоды для всех городов")
	return nil
}

func updateSingleCityWeather(s *weatherService, city model.City, errCh chan<- error) {
	s.log.Info("Начало обновления погоды для города", "city", city.Name)
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
	s.log.Info("Успешное обновление погоды для города", "city", city.Name)
	errCh <- nil
}

func (s *weatherService) GetWeather(city, date string) (model.ResponseFullWeather, error) {
	s.log.Info("Запрос на получение погоды", "city", city, "date", date)
	weather, err := s.wRepo.GetWeatherByCityAndDate(city, date)
	if err != nil {
		s.log.Error("Ошибка при получении погоды", "city", city, "date", date, "error", err)
		return weather, err
	}
	s.log.Info("Погода успешно получена", "city", city, "date", date)
	return weather, nil
}

func (s *weatherService) GetForecast(city string) (model.ResponseShortWeatherByCity, error) {
	s.log.Info("Запрос на получение прогноза погоды", "city", city)
	forecast, err := s.wRepo.GetForecastByCity(city)
	if err != nil {
		s.log.Error("Ошибка при получении прогноза погоды", "city", city, "error", err)
		return forecast, err
	}
	s.log.Info("Прогноз погоды успешно получен", "city", city)
	return forecast, nil
}

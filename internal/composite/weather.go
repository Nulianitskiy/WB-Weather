package composite

import (
	"github.com/jmoiron/sqlx"
	"wb-weather/internal/controller"
	"wb-weather/internal/external"
	"wb-weather/internal/repository"
	"wb-weather/internal/service"
	"wb-weather/pkg/logger"
)

type WeatherComposite struct {
	CityRepository    repository.CityRepo
	WeatherRepository repository.WeatherRepo
	External          external.WeatherExternal
	Service           service.WeatherService
	Controller        controller.WeatherController
}

func NewWeatherComposite(db *sqlx.DB, log logger.Logger, key string) (*WeatherComposite, error) {
	cityRepository := repository.NewCityRepo(db, log)
	weatherRepository := repository.NewWeatherRepo(db, log)
	weatherExternal := external.NewWeatherExternal(key, log)
	weatherService := service.NewWeatherService(weatherRepository, cityRepository, weatherExternal, log)
	weatherController := controller.NewWeatherController(weatherService, log)

	return &WeatherComposite{
		CityRepository:    cityRepository,
		WeatherRepository: weatherRepository,
		External:          weatherExternal,
		Service:           weatherService,
		Controller:        weatherController,
	}, nil
}

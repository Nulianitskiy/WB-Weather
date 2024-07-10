package composite

import (
	"github.com/jmoiron/sqlx"
	"wb-weather/internal/controller"
	"wb-weather/internal/repository"
	"wb-weather/internal/service"
	"wb-weather/pkg/logger"
)

type CityComposite struct {
	Repository repository.CityRepo
	Service    service.CityService
	Controller controller.CityController
}

func NewCityComposite(db *sqlx.DB, log logger.Logger) (*CityComposite, error) {
	cityRepository := repository.NewCityRepo(db, log)
	cityService := service.NewCityService(cityRepository, log)
	cityController := controller.NewCityController(cityService, log)

	return &CityComposite{
		Repository: cityRepository,
		Service:    cityService,
		Controller: cityController}, nil
}

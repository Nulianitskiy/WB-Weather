package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wb-weather/internal/model"
	"wb-weather/internal/service"
	"wb-weather/pkg/logger"
	"wb-weather/pkg/utils"
)

type CityController interface {
	AddCity(ctx *gin.Context)
	GetAllCity(ctx *gin.Context)
}

type cityController struct {
	cityService service.CityService
	log         logger.Logger
}

func NewCityController(cityService service.CityService, log logger.Logger) CityController {
	return &cityController{cityService: cityService, log: log}
}

// AddCity
//
//	@Summary		Добавление нового города
//	@Tags			API дополнительных запросов
//	@Description	Добавление нового города
//	@Accept			json
//	@Produce		json
//
//	@Param			name		query		string	true	"Название города"	example(Москва)
//	@Param			country		query		string	true	"Страна"			example(Россия)
//	@Param			latitude	query		string	true	"Широта"			example(55.7558)
//	@Param			longitude	query		string	true	"Долгота"			example(37.6176)
//
//	@Success		200			{array}		model.City
//	@Failure		400,404		{object}	utils.ErrorResponse
//	@Failure		500			{object}	utils.ErrorResponse
//	@Failure		default		{object}	utils.ErrorResponse
//	@Router			/city [post]
func (cc *cityController) AddCity(ctx *gin.Context) {
	cc.log.Info("Запрос на добавление города")
	var c model.City
	if err := ctx.ShouldBindJSON(&c); err != nil {
		cc.log.Error("Ошибка парсинга параметров города", "error", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
		return
	}

	c, err := cc.cityService.AddCity(c)
	if err != nil {
		cc.log.Error("Ошибка при добавлении города в сервисе", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	cc.log.Info("Город успешно добавлен", "id", c.Id)
	ctx.JSON(http.StatusOK, c)
}

// GetAllCity
//
//	@Summary		Получить все города
//	@Tags			Основное API
//	@Description	Получение списка городов, для которых есть прогноз погоды
//	@Accept			json
//	@Produce		json
//
//	@Success		200		{object}	model.CityList
//	@Failure		400,404	{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Failure		default	{object}	utils.ErrorResponse
//	@Router			/city [get]
func (cc *cityController) GetAllCity(ctx *gin.Context) {
	cc.log.Info("Запрос на получение всех городов с предсказанием погоды")
	cities, err := cc.cityService.GetCities()
	if err != nil {
		cc.log.Error("Ошибка при получении списка городов", "error", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	cc.log.Info("Успешное получение списка городов", "count", len(cities.Cities))
	ctx.JSON(http.StatusOK, cities)
}

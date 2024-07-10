package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
}

func NewCityController(cityService service.CityService) CityController {
	return &cityController{cityService: cityService}
}

func (cc *cityController) AddCity(ctx *gin.Context) {
	var c model.City
	if err := ctx.ShouldBindJSON(&c); err != nil {
		logger.Error("Ошибка парсинга параметров города", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
		return
	}

	c, _ = cc.cityService.AddCity(c)
	logger.Info("Город успешно добавлен")
}

func (cc *cityController) GetAllCity(ctx *gin.Context) {
	c, err := cc.cityService.GetCities()
	if err != nil {
		logger.Error("Ошибка при получении списка городов", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, c)
}

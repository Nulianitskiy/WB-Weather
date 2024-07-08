package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"wb-weather/internal/model"
	"wb-weather/internal/service"
	"wb-weather/pkg/logger"
	"wb-weather/pkg/utils"
)

type WeatherController interface {
	GetWeather(ctx *gin.Context)
	GetForecast(ctx *gin.Context)
}

type weatherController struct {
	weatherService service.WeatherService
}

func NewWeatherController(controller service.WeatherService) WeatherController {
	return &weatherController{weatherService: controller}
}

func (w *weatherController) GetWeather(ctx *gin.Context) {
	weatherId := ctx.Param("id")
	wId, _ := strconv.Atoi(weatherId)

	var weather model.JSONBWeather

	weather, err := w.weatherService.GetWeatherById(wId)
	if err != nil {
		logger.Error("Ошибка при получении списка городов", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, weather)
}

func (w *weatherController) GetForecast(ctx *gin.Context) {
	cityId := ctx.Param("id")
	cId, _ := strconv.Atoi(cityId)

	var forecast []model.ResponseWeather

	forecast, err := w.weatherService.GetForecastByCity(cId)
	if err != nil {
		logger.Error("Ошибка при получении списка городов", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, forecast)
}

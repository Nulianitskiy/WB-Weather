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
	cityName := ctx.Query("city")
	if cityName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}
	date := ctx.Query("date")
	if date == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required"})
		return
	}

	var weather model.ResponseFullWeather
	weather, err := w.weatherService.GetWeather(cityName, date)
	if err != nil {
		logger.Error("Ошибка при получении полной погоды", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, weather)
}

func (w *weatherController) GetForecast(ctx *gin.Context) {
	cityName := ctx.Query("city")
	if cityName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	forecast, err := w.weatherService.GetForecast(cityName)
	if err != nil {
		logger.Error("Ошибка при получении информации по городу", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, forecast)
}

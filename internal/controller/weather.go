package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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
	log            logger.Logger
}

func NewWeatherController(weatherService service.WeatherService, log logger.Logger) WeatherController {
	return &weatherController{weatherService: weatherService, log: log}
}

func (w *weatherController) GetWeather(ctx *gin.Context) {
	cityName := ctx.Query("city")
	if cityName == "" {
		w.log.Error("Не указан параметр city")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}
	date := ctx.Query("date")
	if date == "" {
		w.log.Error("Не указан параметр date")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required"})
		return
	}

	w.log.Info("Запрос на получение погоды", "city", cityName, "date", date)
	weather, err := w.weatherService.GetWeather(cityName, date)
	if err != nil {
		w.log.Error("Ошибка при получении полной погоды", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	w.log.Info("Погода успешно получена", "city", cityName, "date", date)
	ctx.JSON(http.StatusOK, weather)
}

func (w *weatherController) GetForecast(ctx *gin.Context) {
	cityName := ctx.Query("city")
	if cityName == "" {
		w.log.Error("Не указан параметр city")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	w.log.Info("Запрос на получение прогноза погоды", "city", cityName)
	forecast, err := w.weatherService.GetForecast(cityName)
	if err != nil {
		w.log.Error("Ошибка при получении информации по городу", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	w.log.Info("Прогноз погоды успешно получен", "city", cityName)
	ctx.JSON(http.StatusOK, forecast)
}

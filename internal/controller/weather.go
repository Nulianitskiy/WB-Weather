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

// GetWeather
//
//	@Summary		Получить детальную информацию о погоде
//	@Tags			Основное API
//	@Description	Получение детальной информации о погоде для конкретного города и даты
//	@Accept			json
//	@Produce		json
//
//	@Param			city	query		string	true	"Название города"			example(Москва)
//	@Param			date	query		string	true	"Конкретная дата и время"	example(2024-07-10 18:00:00)
//
//	@Success		200		{object}	model.ResponseFullWeather
//	@Failure		400,404	{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Failure		default	{object}	utils.ErrorResponse
//	@Router			/weather [get]
func (w *weatherController) GetWeather(ctx *gin.Context) {
	cityName := ctx.Query("city")
	if cityName == "" {
		w.log.Error("Не указан параметр city")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Не указан параметр city"})
		return
	}
	date := ctx.Query("date")
	if date == "" {
		w.log.Error("Не указан параметр date")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Не указан параметр date"})
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

// GetForecast
//
//	@Summary		Получить предсказание по городу
//	@Tags			Основное API
//	@Description	Получение короткой информации о данных погоды для конкретного города
//	@Accept			json
//	@Produce		json
//
//	@Param			city	query		string	true	"Название города"	example(Москва)
//
//	@Success		200		{object}	model.ResponseShortWeatherByCity
//	@Failure		400,404	{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Failure		default	{object}	utils.ErrorResponse
//	@Router			/forecast [get]
func (w *weatherController) GetForecast(ctx *gin.Context) {
	cityName := ctx.Query("city")
	if cityName == "" {
		w.log.Error("Не указан параметр city")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Не указан параметр city"})
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

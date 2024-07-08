package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"wb-weather/internal/model"
	"wb-weather/internal/repository"
	"wb-weather/pkg/logger"
	"wb-weather/pkg/utils"
)

func GetWeather(ctx *gin.Context) {
	weatherId := ctx.Param("id")
	wId, _ := strconv.Atoi(weatherId)

	var weather model.JSONBWeather

	db, err := repository.GetInstance()
	if err != nil {
		logger.Error("Ошибка получения экземпляра базы данных", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	weather, err = db.GetWeather(wId)
	if err != nil {
		logger.Error("Ошибка при получении списка городов", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, weather)
}

func GetForecast(ctx *gin.Context) {
	cityId := ctx.Param("id")
	cId, _ := strconv.Atoi(cityId)

	var forecast []model.ResponseWeather

	db, err := repository.GetInstance()
	if err != nil {
		logger.Error("Ошибка получения экземпляра базы данных", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	forecast, err = db.GetForecast(cId)
	if err != nil {
		logger.Error("Ошибка при получении списка городов", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, forecast)
}

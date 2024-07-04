package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"wb-weather/internal/model"
	"wb-weather/internal/repository"
	"wb-weather/pkg/logger"
	"wb-weather/pkg/utils"
)

func AddCity(ctx *gin.Context) {
	var c model.City
	if err := ctx.ShouldBindJSON(&c); err != nil {
		logger.Error("Ошибка парсинга параметров города", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
		return
	}

	db, err := repository.GetInstance()
	if err != nil {
		logger.Error("Ошибка получения экземпляра базы данных", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	c, err = db.AddCity(c)
	if err != nil {
		logger.Error("Ошибка при добавлении города", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, c)
	logger.Info("Город успешно добавлен")
}

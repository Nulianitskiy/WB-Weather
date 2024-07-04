package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"wb-weather/internal/repository"
	"wb-weather/internal/routes"
	"wb-weather/pkg/logger"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Ошибка загрузки переменных окружения")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := repository.GetInstance()
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}
	defer db.Close()

	router := gin.Default()

	routes.SetupRoutes(router)

	logger.Info("Запуск сервера на порту", zap.String("port", port))
	err = router.Run(":" + port)
	if err != nil {
		logger.Fatal("Ошибка запуска сервиса", zap.Error(err))
	}

}

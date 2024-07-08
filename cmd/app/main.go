package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"time"
	"wb-weather/internal/controller"
	"wb-weather/internal/external"
	"wb-weather/internal/repository"
	"wb-weather/internal/service"
	"wb-weather/pkg/client/postgresql"
	"wb-weather/pkg/logger"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Ошибка загрузки переменных окружения")
	}
}

// TODO Разгрузить main
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	key := os.Getenv("OWM_KEY")

	database, err := postgresql.NewDatabase()
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}
	defer database.Close()

	router := gin.Default()

	cityRepo := repository.NewCityRepo(database)
	cityService := service.NewCityService(cityRepo)
	cityController := controller.NewCityController(cityService)

	weatherRepo := repository.NewWeatherRepo(database)
	weatherExt := external.NewWeatherExternal(key)
	weatherService := service.NewWeatherService(weatherRepo, cityRepo, weatherExt)
	weatherController := controller.NewWeatherController(weatherService)

	router.POST("/city", cityController.AddCity)
	router.GET("/city", cityController.GetAllCity)
	router.GET("/weather/:id", weatherController.GetWeather)
	router.GET("/forecast/:id", weatherController.GetForecast)

	logger.Info("Запуск сервера на порту", zap.String("port", port))

	ticker := time.NewTicker(3 * time.Hour)
	go func() {
		weatherService.UpdateCityWeather()
		for {
			select {
			case <-ticker.C:
				weatherService.UpdateCityWeather()
			}
		}
	}()

	err = router.Run(":" + port)
	if err != nil {
		logger.Fatal("Ошибка запуска сервиса", zap.Error(err))
	}

}

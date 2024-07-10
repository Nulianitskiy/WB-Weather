package main

import (
	"fmt"
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

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	database, err := postgresql.NewDatabase(dbUser, dbPassword, dbHost, dbPort, dbName)
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
	router.GET("/weather", weatherController.GetWeather)
	router.GET("/forecast", weatherController.GetForecast)

	logger.Info("Запуск сервера на порту", zap.String("port", port))

	ticker := time.NewTicker(3 * time.Hour)
	go func() {
		err := weatherService.UpdateCityWeather()
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			select {
			case <-ticker.C:
				err := weatherService.UpdateCityWeather()
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}()

	err = router.Run(":" + port)
	if err != nil {
		logger.Fatal("Ошибка запуска сервиса", zap.Error(err))
	}

}

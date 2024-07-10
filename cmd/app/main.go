package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "wb-weather/docs"
	"wb-weather/internal/composite"
	"wb-weather/pkg/client/postgresql"
	"wb-weather/pkg/logger"
)

//	@title			Wildberries Weather
//	@version		1.0
//	@description	Тестовое задание для команды портала продавцов..

//	@contact.name	Никита Ульяницкий
//	@contact.url	https://t.me/Nulianitskiy

// @host		localhost:8080
// @BasePath	/
func main() {
	zapLogger := logger.NewZapLogger()

	err := godotenv.Load(".env")
	if err != nil {
		zapLogger.Fatal("Ошибка загрузки переменных окружения", zap.Error(err))
	}

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
		zapLogger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}
	defer database.Close()

	router := gin.Default()

	cityComposite, err := composite.NewCityComposite(database, zapLogger)
	if err != nil {
		zapLogger.Fatal("Ошибка инициализации CityComposite", zap.Error(err))
	}
	weatherComposite, err := composite.NewWeatherComposite(database, zapLogger, key)
	if err != nil {
		zapLogger.Fatal("Ошибка инициализации WeatherComposite", zap.Error(err))
	}

	router.POST("/city", cityComposite.Controller.AddCity)
	router.GET("/city", cityComposite.Controller.GetAllCity)
	router.GET("/weather", weatherComposite.Controller.GetWeather)
	router.GET("/forecast", weatherComposite.Controller.GetForecast)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	zapLogger.Info("Запуск сервера на порту", zap.String("port", port))

	ticker := time.NewTicker(3 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := weatherComposite.Service.UpdateCityWeather()
				if err != nil {
					zapLogger.Error("Ошибка при обновлении погоды", zap.Error(err))
				} else {
					zapLogger.Info("Погода успешно обновлена по расписанию")
				}
			}
		}
	}()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zapLogger.Fatal("Ошибка прослушивания: ", zap.Error(err))
		}
	}()

	// Обработка завершения приложения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	zapLogger.Info("Получен сигнал завершения приложения")

	// Остановка сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Ошибка при остановке сервера", zap.Error(err))
	}
	select {
	case <-ctx.Done():
		zapLogger.Info("Таймаут на 5 секунд.")
	}
	zapLogger.Info("Сервер успешно остановлен")
}

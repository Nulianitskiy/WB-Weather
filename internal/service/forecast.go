package service

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"wb-weather/internal/model"
	"wb-weather/internal/repository"
	"wb-weather/pkg/logger"
)

func InitForecast() {
	db, err := repository.GetInstance()
	if err != nil {
		logger.Error("Ошибка получения экземпляра базы данных", zap.Error(err))
		return
	}

	cities, err := db.GetAllCity()
	if err != nil {
		logger.Error("Ошибка получения городов из бд", zap.Error(err))
		return
	}
	for _, city := range cities {
		CallForecast(city)
	}
}

func CallForecast(c model.City) {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Ошибка загрузки переменных окружения")
	}
	apiKey := os.Getenv("OWM_KEY")
	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%s&lon=%s&units=metric&appid=%s", c.Latitude, c.Longitude, apiKey)

	resp, err := http.Get(apiUrl)
	if err != nil {
		logger.Error("Ошибка получения сообщения с API",
			zap.Error(err),
		)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Ошибка при чтении тела ответа",
			zap.Error(err),
		)
		return
	}

	var weatherData model.WeatherData
	err = json.Unmarshal([]byte(body), &weatherData)
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := repository.GetInstance()
	if err != nil {
		logger.Error("Ошибка получения экземпляра базы данных", zap.Error(err))
		return
	}

	err = db.WeatherUpdateBrute(c, weatherData)
	if err != nil {
		logger.Error("Ошибка добавления данных в базу", zap.Error(err))
		return
	}
}

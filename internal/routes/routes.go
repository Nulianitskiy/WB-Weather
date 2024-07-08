package routes

import (
	"github.com/gin-gonic/gin"
	"wb-weather/internal/controller"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/city", controller.AddCity)
	r.GET("/city", controller.GetAllCity)
	r.GET("/weather/:id", controller.GetWeather)
	r.GET("/forecast/:id", controller.GetForecast)
}

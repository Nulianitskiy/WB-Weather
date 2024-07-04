package routes

import (
	"github.com/gin-gonic/gin"
	"wb-weather/internal/controller"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/city", controller.AddCity)
}

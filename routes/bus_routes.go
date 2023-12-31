package routes

import (
	"github.com/0Mayank/himways/controllers"
	"github.com/gin-gonic/gin"
)

func BusRoutes(router *gin.Engine) {
	router.POST("/bus/", controllers.CreateBus())
	router.GET("/bus/search", controllers.SearchBus())
	router.PUT("/bus/:busId", controllers.EditBus())
	router.GET("/bus/:busId", controllers.GetBus())
	router.DELETE("/bus/:busId", controllers.DeleteBus())
}

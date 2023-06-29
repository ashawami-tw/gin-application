package app

import (
	"github.com/gin-gonic/gin"
)

var route *gin.Engine

func setupRoutes(wiring Wiring) {
	route = gin.Default()
	route.GET("/health", wiring.HealthHandler.HealthHandler)
	route.POST("/create", wiring.CreateUserHandler.CreateUser)
	route.POST("/username", wiring.AddUserNameHandler.AddUserName)
}

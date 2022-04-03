package routes

import (
	"angular-go-jwt/handlers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incoming_routes *gin.Engine){
	incoming_routes.POST("/signup", handlers.SignupHandler)
	incoming_routes.POST("/signin", handlers.SiginHandler)
	incoming_routes.GET("/home", handlers.HomeHandler)
	incoming_routes.GET("/refresh", handlers.RefreshHandler)
	incoming_routes.GET("/logout", handlers.LogoutHandler)

}
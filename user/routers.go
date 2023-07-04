package user

import "github.com/gin-gonic/gin"

func UserRouter(router *gin.RouterGroup) {
	router.POST("/", RegisterUserHandler)
	router.POST("/login", LoginUserHandler)
	
	router.Use(AuthMiddleware(true))
	router.GET("/", GetUserHandler)
}
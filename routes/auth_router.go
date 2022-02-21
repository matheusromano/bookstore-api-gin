package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {
	router.POST("/users/signup", controllers.SignUp())
	router.POST("/users/login", controllers.Login())
}

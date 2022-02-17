package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoute(router *gin.Engine) {
	router.POST("/books", controllers.CreateBook())
	router.GET("/books/:bookId", controllers.GetABook())
	router.PUT("/books/:bookId", controllers.EditABook())
	router.DELETE("books/:bookId", controllers.DeleteABook())
	router.GET("books", controllers.GetAllBooks())
}

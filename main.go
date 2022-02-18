package main

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// run database
	configs.ConnectDB()

	/*router.GET("/", func(c *gin.Context) {
	        c.JSON(200, gin.H{
	                "data": "Gin-gonic & mongoDB already running!",
	        })
	})*/

	//routes
	routes.BookRoute(router)
	routes.SigninRoute(router)

	router.Run("localhost:6000")
}

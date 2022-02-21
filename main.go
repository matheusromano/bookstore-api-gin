package main

import (
	"gin-mongo-api/database"
	"gin-mongo-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	router := gin.Default()

	// run database
	database.ConnectDB()

	//routes
	routes.BookRoute(router)
	//routes.SigninRoute(router)
	routes.AuthRoute(router)

	port := viper.GetString("PORT")

	router.Run(port)
}

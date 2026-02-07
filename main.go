package main

import (
	"pos-service/config"
	"pos-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run("localhost:8080")
}

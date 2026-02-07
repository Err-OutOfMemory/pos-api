package main

import (
	"pos-service/config"
	"pos-service/models"
	"pos-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	
	config.Db.AutoMigrate(
		&models.Employee{},
		&models.User{},
		&models.Category{},
		&models.OrderType{},
		&models.Product{},
		&models.Order{},
		&models.OrderDetail{},
	)

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run("localhost:8080")
}

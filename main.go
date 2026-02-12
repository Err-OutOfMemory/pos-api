package main

import (
	"pos-service/config"
	"pos-service/models"
	"pos-service/routes"
	"time"
	"os"

	"github.com/gin-contrib/cors"
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

	    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{os.Getenv("APP_ADDRESS")},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))
	
	routes.SetupRoutes(r)

	r.Run("localhost:8080")
}

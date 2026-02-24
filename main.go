package main

import (
	"os"
	"pos-service/config"
	"pos-service/models"
	"pos-service/routes"
	"time"

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

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("APP_ADDRESS")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")

	routes.SetupRoutes(r)

	r.Run("localhost:8080")
}

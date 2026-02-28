package main

import (
    "os"
    "pos-service/config"
    "pos-service/models"
    "pos-service/routes"
    "strings"
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

    r.SetTrustedProxies([]string{"172.16.0.0/12", "192.168.0.0/16"})

    r.Use(cors.New(cors.Config{
        AllowOriginFunc: func(origin string) bool {
            return strings.HasPrefix(origin, "http://192.168.") ||
                strings.HasPrefix(origin, "http://localhost") ||
                strings.HasPrefix(origin, "http://127.0.0.1")
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    r.Static("/uploads", "./uploads")

    routes.SetupRoutes(r)

    r.Run(":" + os.Getenv("PORT"))
}

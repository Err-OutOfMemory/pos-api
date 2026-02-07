package routes

import (
	"github.com/gin-gonic/gin"
	"pos-service/controllers"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("", controllers.GetAllEmployees)
			users.GET("/:id", controllers.GetEmployeeByID)
			users.POST("", controllers.CreateEmployee)
			users.PUT("/:id", controllers.UpdateEmployee)
			users.DELETE("/:id", controllers.DeleteEmployee)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", controllers.GetCategories)
		}
	}
}

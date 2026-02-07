package routes

import (
	"github.com/gin-gonic/gin"
	"pos-service/controllers"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		employees := api.Group("/employees")
		{
			employees.GET("", controllers.GetAllEmployees)
			employees.GET("/:id", controllers.GetEmployeeByID)
			employees.POST("", controllers.CreateEmployee)
			employees.PUT("/:id", controllers.UpdateEmployee)
			employees.DELETE("/:id", controllers.DeleteEmployee)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", controllers.GetCategories)
		}
	}
}

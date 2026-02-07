package routes

import (
	"github.com/gin-gonic/gin"
	"pos-service/controllers"
	"pos-service/middleware"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		//Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/check_user", controllers.CheckUser)
			auth.POST("/set_pin", controllers.SetupPin)
			auth.POST("/login", controllers.Login)
		}

		//Protected routes
		protected := api.Group("", middleware.AuthMiddleware())
		{
			employees := protected.Group("/employees", middleware.AuthorizeRole("admin"))
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
				categories.GET("/:id", controllers.GetCategoryByID)
				categories.POST("", controllers.CreateCategory)
				categories.PUT("/:id", controllers.UpdateCategory)
				categories.DELETE("/:id", controllers.DeleteCategory)
			}

			products := api.Group("/products")
			{
				products.GET("", controllers.GetProducts)
				products.GET("/:id", controllers.GetProductByID)
				products.POST("", controllers.CreateProduct)
				products.PUT("/:id", controllers.UpdateProduct)
				products.DELETE("/:id", controllers.DeleteProduct)
			}

			orderTypes := api.Group("/order_types")
			{
				orderTypes.GET("", controllers.GetOrderTypes)
				orderTypes.GET("/:id", controllers.GetOrderTypeByID)
				orderTypes.POST("", controllers.CreateOrderType)
				orderTypes.PUT("/:id", controllers.UpdateOrderType)
				orderTypes.DELETE("/:id", controllers.DeleteOrderType)
			}
		}
	}
}

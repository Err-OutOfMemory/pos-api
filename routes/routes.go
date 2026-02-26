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
			auth.GET("/profile", middleware.AuthMiddleware(), controllers.GetProfile)
		}

		//Protected routes
		protected := api.Group("", middleware.AuthMiddleware())
		{
			protected.POST("/upload", controllers.UploadFile)

			employees := protected.Group("/employees", middleware.AuthorizeRole("admin"))
			{
				employees.GET("", controllers.GetAllEmployees)
				employees.GET("/:id", controllers.GetEmployeeByID)
				employees.POST("", controllers.CreateEmployee)
				employees.PUT("/:id", controllers.UpdateEmployee)
				employees.DELETE("/:id", controllers.DeleteEmployee)
			}

			categories := protected.Group("/categories")
			{
				categories.GET("", controllers.GetCategories)
				categories.GET("/:id", controllers.GetCategoryByID)
				categories.POST("", controllers.CreateCategory)
				categories.PUT("/:id", controllers.UpdateCategory)
				categories.DELETE("/:id", controllers.DeleteCategory)
			}

			products := protected.Group("/products")
			{
				products.GET("", controllers.GetProducts)
				products.GET("/:id", controllers.GetProductByID)
				products.POST("", controllers.CreateProduct)
				products.PUT("/:id", controllers.UpdateProduct)
				products.DELETE("/:id", controllers.DeleteProduct)
			}

			orderTypes := protected.Group("/order_types")
			{
				orderTypes.GET("", controllers.GetOrderTypes)
				orderTypes.GET("/:id", controllers.GetOrderTypeByID)
				orderTypes.POST("", controllers.CreateOrderType)
				orderTypes.PUT("/:id", controllers.UpdateOrderType)
				orderTypes.DELETE("/:id", controllers.DeleteOrderType)
			}

			order := protected.Group("/orders")
			{
				order.GET("", controllers.GetOrders)
				order.GET("/:id", controllers.GetOrderByID)
				order.POST("", controllers.CreateOrder)
				order.PUT("/:id", controllers.UpdateOrder)
				order.PATCH("/:id/status", controllers.UpdateOrderStatus)
				order.DELETE("/:id", controllers.DeleteOrder)
			}
		}
	}
}

package routes

import (
	"manajemen-karyawan-api/controller"
	"manajemen-karyawan-api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine) {
	// Apply CORS globally
	r.Use(middleware.CORSMiddleware())

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Group: /api
	api := r.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", controller.Login)
			auth.POST("/logout", controller.Login)
			auth.GET("/me", middleware.AuthMiddleware(), controller.GetMe)
		}

		// Protected routes (cookie-based JWT)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())

		// Employee routes
		employee := protected.Group("/employee")
		{
			employee.POST("/GetData", controller.GetAllEmployees)
			employee.GET("/:employee_id", controller.GetEmployeeByID)
			employee.POST("", controller.CreateEmployee)
			employee.PUT("/:id", controller.UpdateEmployee)
			employee.DELETE("/:id", controller.DeleteEmployee)
		}

		// Departement routes
		departement := protected.Group("/departement")
		{
			departement.POST("/GetData", controller.GetAllDepartements)
			departement.GET("/:id", controller.GetDepartementByID)
			departement.POST("", controller.CreateDepartement)
			departement.PUT("/:id", controller.UpdateDepartement)
			departement.DELETE("/:id", controller.DeleteDepartement)
		}

		//  Attendance routes
		attendance := protected.Group("/attendance")
		{
			attendance.POST("", controller.ClockHandler)
			attendance.GET("/today", controller.GetTodayAttendance)
			attendance.POST("/logs", controller.GetAttendanceLogs)
			attendance.POST("/GetData", controller.GetAllAttendanceLogs)

		}
	}
}

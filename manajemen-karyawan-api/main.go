package main

import (
	"log"
	"manajemen-karyawan-api/config"
	"manajemen-karyawan-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "manajemen-karyawan-api/docs"
)

// @title Employee Management API
// @version 1.0
// @description API for managing employee data
// @termsOfService https://example.com/terms/

// @contact.name Backend Team
// @contact.email backend@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_token

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to system environment variables")
	}

	// ✅ Load all configuration from environment
	config.InitConfig()

	// ✅ Initialize database connection
	config.InitDB()
	defer config.DB.Close()

	// ✅ Initialize Gin router
	r := gin.Default()

	// ✅ Register all routes
	routes.RegisterRoutes(r)

	// ✅ Get port from environment or use default
	port := config.AppPort
	if port == "" {
		port = "8080"
	}

	// ✅ Start the server
	log.Printf("Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

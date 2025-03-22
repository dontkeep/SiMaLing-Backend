package main

import (
	"log"

	"github.com/dontkeep/simaling-backend/controllers"
	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.DatabaseConnection()
	if err := controllers.CreateDefaultRoles(); err != nil {
		log.Fatalf("Failed to create default roles: %v", err)
	}
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/login", controllers.Login)

	authorized := r.Group("/api")

	authorized.Use(controllers.Authenticate)
	{
		// User management
		authorized.GET("/users", controllers.GetAllUsers)
		authorized.GET("/users/:id", controllers.GetUser)
		authorized.POST("/users", controllers.CreateUser)
		authorized.PUT("/users/:id", controllers.UpdateUser)
		authorized.DELETE("/users/:id", controllers.DeleteUser)

		// Role management
		authorized.GET("/roles", controllers.GetAllRoles)

		// Funds management
		authorized.GET("/funds", controllers.GetFunds)
		authorized.GET("/funds/:id", controllers.GetFundsById)
		authorized.GET("/funds-by-user", controllers.GetFundsByUser)
		authorized.POST("/funds", controllers.CreateFunds)
		authorized.PUT("/funds/:id", controllers.UpdateFunds)
		authorized.DELETE("/funds/:id", controllers.DeleteFunds)

		// Security records management
		authorized.GET("/security-records", controllers.GetAllSecurityRecord)        // Admin-only
		authorized.POST("/security-records", controllers.CreateSecurityRecord)       // Security-only
		authorized.DELETE("/security-records/:id", controllers.DeleteSecurityRecord) // Admin-only

		// Logout
		authorized.POST("/logout", controllers.Logout)
	}

	r.Run()
}

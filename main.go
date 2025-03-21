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
		authorized.GET("/users", controllers.GetAllUsers)
		authorized.GET("/users/:id", controllers.GetUser)
		authorized.GET("/roles", controllers.GetAllRoles)
		authorized.GET("/funds", controllers.GetFunds)
		authorized.POST("/users", controllers.CreateUser)
		authorized.POST("/funds", controllers.CreateFunds)
		authorized.POST("/logout", controllers.Logout)
	}

	r.Run()
}

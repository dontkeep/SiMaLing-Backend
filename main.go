package main

import (
	"log"

	"github.com/dontkeep/simaling-backend/controllers"
	"github.com/dontkeep/simaling-backend/initializers"
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
	r.GET("/roles", controllers.GetAllRoles)
	r.Run()
}

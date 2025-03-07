package main

import (
	"github.com/dontkeep/simaling-backend/controllers"
	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.DatabaseConnection()
}

func main() {
	r := gin.Default()
	r.POST("/post", controllers.PostCreate)
	r.GET("/post", controllers.GetAllPost)
	r.GET("/post/:id", controllers.GetPost)
	r.PUT("/post/:id", controllers.UpdatePost)
	r.DELETE("/post/:id", controllers.DeletePost)
	r.Run()
}

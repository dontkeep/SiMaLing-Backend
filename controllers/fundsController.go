package controllers

import (
	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
)

func GetFunds(c *gin.Context) {
	var funds []models.Funds

	result := initializers.DB.Find(&funds)

	if result != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get funds",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": funds,
	})
}

// func CreateFunds(c * gin.Context) {
// 	var body struct {
// 		Amount int

// 	}
// }

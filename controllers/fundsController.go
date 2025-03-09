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

func CreateFunds(c *gin.Context) {
	var body struct {
		Amount      int
		Is_Income   bool
		Description string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{
			"message": "User not authenticated",
		})
		return
	}
	uid, ok := userId.(uint)
	if !ok {
		c.JSON(400, gin.H{
			"message": "Invalid user id",
		})
	}

	funds := models.Funds{
		User_Id:     uid,
		Amount:      float64(body.Amount),
		Is_Income:   body.Is_Income,
		Description: body.Description,
	}

	result := initializers.DB.Create(&funds)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to create funds",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Funds created",
		"data":    funds,
	})
}

func DeleteFunds(c *gin.Context) {
	id := c.Param("id")
	var funds models.Funds

	result := initializers.DB.Where("id = ?", id).First(&funds)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get funds",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Funds deleted",
		"data":    funds,
	})
}

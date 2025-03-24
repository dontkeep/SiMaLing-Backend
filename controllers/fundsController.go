package controllers

import (
	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
)

// GetFunds gets all funds
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

// CreateFunds creates a new funds
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
		Status:      "Pending",
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

// DeleteFunds deletes a funds
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

// UpdateFunds updates a funds
func UpdateFunds(c *gin.Context) {
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

	id := c.Param("id")
	var funds models.Funds

	result := initializers.DB.Where("id = ?", id).First(&funds)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get funds",
		})
		return
	}

	funds.Amount = float64(body.Amount)
	funds.Is_Income = body.Is_Income
	funds.Description = body.Description

	result = initializers.DB.Save(&funds)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to update funds",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Funds updated",
		"data":    funds,
	})
}

// GetFundsByUser gets all funds by user
func GetFundsByUser(c *gin.Context) {
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

	var funds []models.Funds

	result := initializers.DB.Where("user_id = ?", uid).Find(&funds)

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

// GetFundsById gets a funds by id'
func GetFundsById(c *gin.Context) {
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
		"data": funds,
	})
}

func AcceptFunds(c *gin.Context) {
	id := c.Param("id")
	var funds models.Funds

	// Retrieve the funds record by ID
	result := initializers.DB.Where("id = ?", id).First(&funds)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get funds",
		})
		return
	}

	// Update the status to "Accepted"
	funds.Status = "Accepted"

	// Save the updated funds record
	result = initializers.DB.Save(&funds)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to accept funds",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Funds accepted",
		"data":    funds,
	})
}

func RejectFunds(c *gin.Context) {
	id := c.Param("id")
	var funds models.Funds

	// Retrieve the funds record by ID
	result := initializers.DB.Where("id = ?", id).First(&funds)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get funds",
		})
		return
	}

	// Update the status to "Rejected"
	funds.Status = "Rejected"

	// Save the updated funds record
	result = initializers.DB.Save(&funds)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to reject funds",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Funds rejected",
		"data":    funds,
	})
}

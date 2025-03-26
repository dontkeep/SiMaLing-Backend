package controllers

import (
	"strconv"

	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
)

// GetFunds gets all funds
func GetFunds(c *gin.Context) {
	// Get query parameters for pagination
	page := c.DefaultQuery("page", "1")    // Default to page 1 if not provided
	limit := c.DefaultQuery("limit", "10") // Default to 10 records per page if not provided

	// Convert query parameters to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(400, gin.H{
			"message": "Invalid page number",
		})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		c.JSON(400, gin.H{
			"message": "Invalid limit number",
		})
		return
	}

	// Calculate the offset
	offset := (pageInt - 1) * limitInt

	// Retrieve paginated funds from the database
	var funds []models.Funds
	result := initializers.DB.Limit(limitInt).Offset(offset).Find(&funds)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get funds",
		})
		return
	}

	// Count the total number of funds
	var total int64
	initializers.DB.Model(&models.Funds{}).Count(&total)

	// Return the paginated response
	c.JSON(200, gin.H{
		"data":       funds,
		"total":      total,
		"page":       pageInt,
		"limit":      limitInt,
		"totalPages": (total + int64(limitInt) - 1) / int64(limitInt), // Calculate total pages
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
	// Get query parameters for pagination
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(400, gin.H{
			"message": "Invalid page number",
		})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		c.JSON(400, gin.H{
			"message": "Invalid limit number",
		})
		return
	}

	offset := (pageInt - 1) * limitInt

	// Retrieve user ID from context
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
			"message": "Invalid user ID",
		})
		return
	}

	// Retrieve paginated funds for the user
	var funds []models.Funds
	result := initializers.DB.Where("user_id = ?", uid).Limit(limitInt).Offset(offset).Find(&funds)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get funds",
		})
		return
	}

	// Count the total number of funds for the user
	var total int64
	initializers.DB.Model(&models.Funds{}).Where("user_id = ?", uid).Count(&total)

	// Return the paginated response
	c.JSON(200, gin.H{
		"data":       funds,
		"total":      total,
		"page":       pageInt,
		"limit":      limitInt,
		"totalPages": (total + int64(limitInt) - 1) / int64(limitInt),
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

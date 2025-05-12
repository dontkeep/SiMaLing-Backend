package controllers

import (
	"fmt"
	"strconv"

	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
)

// if role is not admin, return 403
// if role is admin, return all security record

func GetScurityRecordToday(c *gin.Context) {
	// function
}

func GetAllSecurityRecord(c *gin.Context) {
	// Check if the role is admin
	role := c.MustGet("role").(string)
	if role != "admin" {
		c.JSON(403, gin.H{
			"message": "Forbidden",
		})
		return
	}

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

	// Retrieve paginated security records
	var securityRecords []models.SecurityRecord
	result := initializers.DB.Limit(limitInt).Offset(offset).Find(&securityRecords)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get security records",
		})
		return
	}

	// Count the total number of security records
	var total int64
	initializers.DB.Model(&models.SecurityRecord{}).Count(&total)

	// Return the paginated response
	c.JSON(200, gin.H{
		"data":       securityRecords,
		"total":      total,
		"page":       pageInt,
		"limit":      limitInt,
		"totalPages": (total + int64(limitInt) - 1) / int64(limitInt),
	})
}

func GetSecurityRecordByPhoneNum(c *gin.Context) {
	// Retrieve user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{
			"message": "User not authenticated",
		})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(400, gin.H{
			"message": "Invalid user ID",
		})
		return
	}

	// Retrieve the user by ID
	var user models.User
	if err := initializers.DB.Where("id = ?", uid).First(&user).Error; err != nil {
		c.JSON(404, gin.H{
			"message": "User not found",
		})
		return
	}

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

	// Retrieve paginated security records by the user's NIK
	var securityRecords []models.SecurityRecord
	result := initializers.DB.Where("security_id = ?", user.ID).Limit(limitInt).Offset(offset).Find(&securityRecords)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "Failed to retrieve security records",
		})
		return
	}

	// Count the total number of security records for the user
	var total int64
	initializers.DB.Model(&models.SecurityRecord{}).Where("security_id = ?", user.ID).Count(&total)

	// Return the paginated response
	c.JSON(200, gin.H{
		"data":       securityRecords,
		"total":      total,
		"page":       pageInt,
		"limit":      limitInt,
		"totalPages": (total + int64(limitInt) - 1) / int64(limitInt),
	})
}

// only security can create security record
func CreateSecurityRecord(c *gin.Context) {
	// Get the role from the context
	role := c.MustGet("role").(string)
	if role != "security" {
		c.JSON(403, gin.H{
			"message": "Forbidden",
		})
		return
	}

	// Parse the request body
	var body struct {
		UserID    uint    `json:"user_id"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		Phone_No  string  `json:"phone_no"`
		Note      string  `json:"note"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Verify if a user with the provided NIK exists
	var user models.User
	if err := initializers.DB.Where("phone_no = ?", body.Phone_No).First(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "No user found with the provided Phone Number",
		})
		return
	}

	// Create the security record
	securityRecord := models.SecurityRecord{
		Security_Id: body.UserID,
		Phone_No:    body.Phone_No,
		Longitude:   fmt.Sprintf("%f", body.Longitude),
		Latitude:    fmt.Sprintf("%f", body.Latitude),
	}

	// Save the record to the database
	if err := initializers.DB.Create(&securityRecord).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to create security record",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Security record created successfully",
		"data":    securityRecord,
	})
}

func DeleteSecurityRecord(c *gin.Context) {
	// Get the role from the context
	role := c.MustGet("role").(string)
	if role != "admin" {
		c.JSON(403, gin.H{
			"message": "Forbidden",
		})
		return
	}

	// Retrieve the security record ID from the URL
	id := c.Param("id")

	// Retrieve the security record by ID
	var securityRecord models.SecurityRecord
	result := initializers.DB.Where("id = ?", id).First(&securityRecord)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get security record",
		})
		return
	}

	// Delete the security record
	result = initializers.DB.Delete(&securityRecord)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to delete security record",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Security record deleted",
	})
}

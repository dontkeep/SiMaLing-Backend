package controllers

import (
	"fmt"

	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
)

// if role is not admin, return 403
// if role is admin, return all security record
func GetAllSecurityRecord(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "admin" {
		c.JSON(403, gin.H{
			"message": "Forbidden",
		})
		return
	}
	var securityRecord []models.SecurityRecord
	initializers.DB.Find(&securityRecord)
	c.JSON(200, securityRecord)
}

func DeleteSecurityRecord(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "admin" {
		c.JSON(403, gin.H{
			"message": "Forbidden",
		})
		return
	}
	id := c.Param("id")
	var securityRecord models.SecurityRecord
	initializers.DB.First(&securityRecord, id)
	initializers.DB.Delete(&securityRecord)
	c.JSON(200, gin.H{
		"message": "Security record deleted",
	})
}

func GetSecurityRecordByNik(c *gin.Context) {
	// receive token from context, get the user id from the token, and get the user nik from the user id, then get the security record by the nik, can be used by all roles
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{
			"message": "User not authenticated",
		})
		return
	}

	// Type assert userID to uint
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

	// Retrieve security records by the user's NIK
	var securityRecords []models.SecurityRecord
	if err := initializers.DB.Where("security_id = ?", user.ID).Find(&securityRecords).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to retrieve security records",
		})
		return
	}

	// Return the security records
	c.JSON(200, gin.H{
		"message": "Security records retrieved successfully",
		"data":    securityRecords,
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
		NIK       string  `json:"nik"`
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
	if err := initializers.DB.Where("nik = ?", body.NIK).First(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "No user found with the provided NIK",
		})
		return
	}

	// Create the security record
	securityRecord := models.SecurityRecord{
		Security_Id: body.UserID,
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

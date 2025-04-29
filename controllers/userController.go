package controllers

import (
	"fmt"
	"strconv"

	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdminAccount() error {
	// Check if the role exists
	var role models.Roles
	if err := initializers.DB.Where("id = ?", 1).First(&role).Error; err != nil {
		return fmt.Errorf("role not found: %v", err)
	}

	// Check if an admin user already exists
	var user models.User
	initializers.DB.Where("role_id = ?", 1).First(&user)
	if user.ID != 0 {
		return nil // Admin account already exists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Create the admin account
	user = models.User{
		Phone_No: "081234567890",
		NIK:      "1234567890123456",
		Password: string(hashedPassword),
		Name:     "Admin",
		Address:  "Jl. Admin",
		Role_Id:  role.ID, // Use the Role ID from the Roles table
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllUsers(c *gin.Context) {
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

	type UserResponse struct {
		ID       uint   `json:"id"`
		Phone_No string `json:"phone_no"`
		NIK      string `json:"nik"`
		Name     string `json:"name"`
		Address  string `json:"address"`
		Role_Id  uint   `json:"role_id"`
	}

	// Retrieve paginated users from the database
	var users []UserResponse
	result := initializers.DB.Model(&models.User{}).Select("id, phone_no, nik, name, address, role_id").Limit(limitInt).Offset(offset).Scan(&users)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get users",
		})
		return
	}

	// Count the total number of users
	var total int64
	initializers.DB.Model(&models.User{}).Count(&total)

	// Return the paginated response
	c.JSON(200, gin.H{
		"data":       users,
		"total":      total,
		"page":       pageInt,
		"limit":      limitInt,
		"totalPages": (total + int64(limitInt) - 1) / int64(limitInt), // Calculate total pages
	})
}

func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	result := initializers.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get user",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": user,
	})
}

func CreateUser(c *gin.Context) {
	var body struct {
		Phone_No      string
		NIK           string
		Password      string
		Name          string
		Address       string
		Role_Id       uint
		FamilyMembers []struct {
			NIK    string `json:"nik"`
			Name   string `json:"name"`
			Status string `json:"status"` // "wife" or "child"
		} `json:"family_members"`
	}

	// Parse the request body
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Create the user
	user := models.User{
		Phone_No: body.Phone_No,
		NIK:      body.NIK,
		Password: body.Password,
		Name:     body.Name,
		Address:  body.Address,
		Role_Id:  body.Role_Id,
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to create user",
		})
		return
	}

	// Create family members if provided
	if len(body.FamilyMembers) > 0 {
		var familyMembers []models.FamilyMembers
		for _, fm := range body.FamilyMembers {
			familyMembers = append(familyMembers, models.FamilyMembers{
				NIK:             fm.NIK,
				Name:            fm.Name,
				Status:          fm.Status,
				HeadOfFamily_Id: user.ID,
			})
		}

		result = initializers.DB.Create(&familyMembers)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"message": "Failed to create family members",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "User and family members created successfully",
		"user":    user,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Check if the user exists
	result := initializers.DB.First(&user, id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"message": "User not found",
		})
		return
	}

	// Parse the request body
	var body struct {
		Phone_No      string
		NIK           string
		Password      string
		Name          string
		Address       string
		Role_Id       uint
		FamilyMembers []struct {
			ID     uint   `json:"id"` // Include ID to identify existing family members
			NIK    string `json:"nik"`
			Name   string `json:"name"`
			Status string `json:"status"` // "wife" or "child"
		} `json:"family_members"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Update the user fields
	user.Phone_No = body.Phone_No
	user.NIK = body.NIK
	user.Password = body.Password
	user.Name = body.Name
	user.Address = body.Address
	user.Role_Id = body.Role_Id

	// Save the updated user
	result = initializers.DB.Save(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to update user",
		})
		return
	}

	// Process family members
	var existingFamilyMembers []models.FamilyMembers
	initializers.DB.Where("head_of_family_id = ?", user.ID).Find(&existingFamilyMembers)

	// Map existing family members by ID for quick lookup
	existingFamilyMap := make(map[uint]models.FamilyMembers)
	for _, fm := range existingFamilyMembers {
		existingFamilyMap[fm.ID] = fm
	}

	// Track IDs of family members in the request
	requestedFamilyIDs := make(map[uint]bool)

	for _, fm := range body.FamilyMembers {
		requestedFamilyIDs[fm.ID] = true

		if fm.ID == 0 {
			// Add new family member
			newFamilyMember := models.FamilyMembers{
				NIK:             fm.NIK,
				Name:            fm.Name,
				Status:          fm.Status,
				HeadOfFamily_Id: user.ID,
			}
			initializers.DB.Create(&newFamilyMember)
		} else {
			// Update existing family member
			if existingFamily, exists := existingFamilyMap[fm.ID]; exists {
				existingFamily.NIK = fm.NIK
				existingFamily.Name = fm.Name
				existingFamily.Status = fm.Status
				initializers.DB.Save(&existingFamily)
			}
		}
	}

	// Delete family members not in the request
	for _, existingFamily := range existingFamilyMembers {
		if !requestedFamilyIDs[existingFamily.ID] {
			initializers.DB.Delete(&existingFamily)
		}
	}

	c.JSON(200, gin.H{
		"message": "User and family members have been updated",
		"user":    user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	result := initializers.DB.Delete(&user, id)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Failed to delete user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User has been deleted",
	})
}

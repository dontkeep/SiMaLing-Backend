package controllers

import (
	"fmt"
	"strconv"
	"time"

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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("thisissecured"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Create the admin account
	user = models.User{
		Phone_No: "082298588849",
		Email:    "donnya238@gmail.com",
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
		Email    string `json:"email"`
		Name     string `json:"name"`
		Address  string `json:"address"`
		Role_Id  uint   `json:"role_id"`
	}

	// Retrieve paginated users from the database
	var users []UserResponse
	result := initializers.DB.Model(&models.User{}).Select("id, email, phone_no, name, address, role_id").Limit(limitInt).Offset(offset).Scan(&users)
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

func getHomeData(c *gin.Context) {
	// Get total users
	var totalUsers int64
	initializers.DB.Model(&models.User{}).Count(&totalUsers)

	// Get month and year from query, default to current month/year
	month := time.Now().Month()
	year := time.Now().Year()
	if m := c.Query("month"); m != "" {
		if mi, err := strconv.Atoi(m); err == nil && mi >= 1 && mi <= 12 {
			month = time.Month(mi)
		}
	}
	if y := c.Query("year"); y != "" {
		if yi, err := strconv.Atoi(y); err == nil && yi > 0 {
			year = yi
		}
	}

	// Calculate start and end of the month
	startTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endTime := startTime.AddDate(0, 1, 0)

	// Get total users added in the month
	var usersAddedThisMonth int64
	initializers.DB.Model(&models.User{}).
		Where("created_at >= ? AND created_at < ?", startTime, endTime).
		Count(&usersAddedThisMonth)

	// Get total income and total expense for the month
	var totalIncome float64
	initializers.DB.Model(&models.Funds{}).
		Select("COALESCE(SUM(amount),0)").
		Where("is_income = ? AND created_at >= ? AND created_at < ?", true, startTime, endTime).
		Row().Scan(&totalIncome)

	var totalExpense float64
	initializers.DB.Model(&models.Funds{}).
		Select("COALESCE(SUM(amount),0)").
		Where("is_income = ? AND created_at >= ? AND created_at < ?", false, startTime, endTime).
		Row().Scan(&totalExpense)

	c.JSON(200, gin.H{
		"total_users":            totalUsers,
		"users_added_this_month": usersAddedThisMonth,
		"total_income":           totalIncome,
		"total_expense":          totalExpense,
		"month":                  int(month),
		"year":                   year,
	})
}

// Exported version of getHomeData for routing
func GetHomeData(c *gin.Context) {
	getHomeData(c)
}

func CreateUser(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(403, gin.H{"message": "Forbidden: Admins only"})
		return
	}

	var body struct {
		Phone_No      string `json:"phone_no"`
		Email         string `json:"email"`
		Password      string `json:"password"`
		Name          string `json:"name"`
		Address       string `json:"address"`
		Role_Id       uint   `json:"role_id"`
		FamilyMembers []struct {
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

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to hash password"})
		return
	}

	// Create the user
	user := models.User{
		Phone_No: body.Phone_No,
		Email:    body.Email,
		Password: string(hashedPassword),
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

	type UserResponse struct {
		ID       uint   `json:"id"`
		Phone_No string `json:"phone_no"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Address  string `json:"address"`
		Role_Id  uint   `json:"role_id"`
	}

	response := UserResponse{
		ID:       user.ID,
		Phone_No: user.Phone_No,
		Email:    user.Email,
		Name:     user.Name,
		Address:  user.Address,
		Role_Id:  user.Role_Id,
	}

	c.JSON(200, gin.H{
		"message": "User and family members created successfully",
		"user":    response,
	})
}

func UpdateUser(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(403, gin.H{"message": "Forbidden: Admins only"})
		return
	}

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
		Email         string
		Password      string
		Name          string
		Address       string
		Role_Id       uint
		FamilyMembers []struct {
			ID     uint   `json:"id"` // Include ID to identify existing family members
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
	user.Email = body.Email
	user.Name = body.Name
	user.Address = body.Address
	user.Role_Id = body.Role_Id

	// Hash the password if it is being updated
	if body.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

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
				Name:            fm.Name,
				Status:          fm.Status,
				HeadOfFamily_Id: user.ID,
			}
			initializers.DB.Create(&newFamilyMember)
		} else {
			// Update existing family member
			if existingFamily, exists := existingFamilyMap[fm.ID]; exists {
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

	// Prepare response struct
	type UserResponse struct {
		ID       uint   `json:"id"`
		Phone_No string `json:"phone_no"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Address  string `json:"address"`
		Role_Id  uint   `json:"role_id"`
	}
	response := UserResponse{
		ID:       user.ID,
		Phone_No: user.Phone_No,
		Email:    user.Email,
		Name:     user.Name,
		Address:  user.Address,
		Role_Id:  user.Role_Id,
	}

	c.JSON(200, gin.H{
		"message": "User and family members have been updated",
		"user":    response,
	})
}

func DeleteUser(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(403, gin.H{"message": "Forbidden: Admins only"})
		return
	}

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

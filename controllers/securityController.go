package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
)

// userVerification checks if the current user is a normal user (role_id == 2)
func userVerification(c *gin.Context) bool {
	userID, exists := c.Get("user_id")
	if !exists {
		return false
	}
	uid, ok := userID.(uint)
	if !ok {
		return false
	}
	var user models.User
	if err := initializers.DB.First(&user, uid).Error; err != nil {
		return false
	}
	return user.Role_Id == 2 // 2 = Normal user
}

// if role is not admin, return 403
// if role is admin, return all security record

func GetScurityRecordToday(c *gin.Context) {
	if !userVerification(c) {
		c.JSON(403, gin.H{"message": "Forbidden: Users only"})
		return
	}

	// Get user id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"message": "User not authenticated"})
		return
	}
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(400, gin.H{"message": "Invalid user ID"})
		return
	}

	// Get today's date (YYYY-MM-DD)
	today := time.Now().Format("2006-01-02")

	// Get user's phone number
	var user models.User
	if err := initializers.DB.Select("phone_no").Where("id = ?", uid).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"message": "User not found"})
		return
	}

	type SecurityRecordResponse struct {
		ID           uint   `json:"id"`
		SecurityId   uint   `json:"security_id"`
		SecurityName string `json:"security_name"`
		Block        string `json:"block"`
		PhoneNo      string `json:"phone_no"`
		Longitude    string `json:"longitude"`
		Latitude     string `json:"latitude"`
	}

	var records []SecurityRecordResponse
	result := initializers.DB.Model(&models.SecurityRecord{}).
		Select("security_records.id, security_records.security_id, users.name as security_name, security_records.block, security_records.phone_no, security_records.longitude, security_records.latitude").
		Joins("left join users on users.id = security_records.security_id").
		Where("security_records.phone_no = ? AND DATE(security_records.created_at) = ?", user.Phone_No, today).
		Scan(&records)
	if result.Error != nil {
		c.JSON(400, gin.H{"message": "Failed to get today's security records"})
		return
	}

	c.JSON(200, gin.H{"data": records})
}

func GetAllSecurityRecord(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(403, gin.H{"message": "Forbidden: Admins only"})
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

	type SecurityRecordResponse struct {
		ID           uint   `json:"id"`
		SecurityId   uint   `json:"security_id"`
		SecurityName string `json:"security_name"`
		Block        string `json:"block"`
		PhoneNo      string `json:"phone_no"`
		Longitude    string `json:"longitude"`
		Latitude     string `json:"latitude"`
	}

	var records []SecurityRecordResponse
	result := initializers.DB.Model(&models.SecurityRecord{}).
		Select("security_records.id, security_records.security_id, users.name as security_name, security_records.block, security_records.phone_no, security_records.longitude, security_records.latitude").
		Joins("left join users on users.id = security_records.security_id").
		Limit(limitInt).Offset(offset).
		Scan(&records)
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
		"data":       records,
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
		Block     string  `json:"block"`
		Phone_No  string  `json:"phone_no"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Get security user id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"message": "User not authenticated"})
		return
	}
	securityID, ok := userID.(uint)
	if !ok {
		c.JSON(400, gin.H{"message": "Invalid user ID"})
		return
	}

	// Get security name
	var securityUser models.User
	if err := initializers.DB.Select("name").Where("id = ?", securityID).First(&securityUser).Error; err != nil {
		c.JSON(400, gin.H{"message": "Security user not found"})
		return
	}

	// Create the security record
	securityRecord := models.SecurityRecord{
		Security_Id: securityID,
		Block:       body.Block,
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

	// Prepare response
	response := struct {
		ID           uint   `json:"id"`
		SecurityId   uint   `json:"security_id"`
		SecurityName string `json:"security_name"`
		Block        string `json:"block"`
		PhoneNo      string `json:"phone_no"`
		Longitude    string `json:"longitude"`
		Latitude     string `json:"latitude"`
	}{
		ID:           securityRecord.ID,
		SecurityId:   securityID,
		SecurityName: securityUser.Name,
		Block:        securityRecord.Block,
		PhoneNo:      securityRecord.Phone_No,
		Longitude:    securityRecord.Longitude,
		Latitude:     securityRecord.Latitude,
	}

	c.JSON(200, gin.H{
		"message": "Security record created successfully",
		"data":    response,
	})
}

func AddSecurityRecord(c *gin.Context) {
	if !isSecurity(c) {
		c.JSON(403, gin.H{"message": "Forbidden: Security only"})
		return
	}

	// Parse the request body (only block, longitude, latitude)
	var body struct {
		Block     string  `json:"block"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	// Get security user id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"message": "User not authenticated"})
		return
	}
	securityID, ok := userID.(uint)
	if !ok {
		c.JSON(400, gin.H{"message": "Invalid user ID"})
		return
	}

	// Get security user (for name and phone_no)
	var securityUser models.User
	if err := initializers.DB.Select("name, phone_no").Where("id = ?", securityID).First(&securityUser).Error; err != nil {
		c.JSON(400, gin.H{"message": "Security user not found"})
		return
	}

	// Create the security record
	securityRecord := models.SecurityRecord{
		Security_Id: securityID,
		Block:       body.Block,
		Phone_No:    securityUser.Phone_No,
		Longitude:   fmt.Sprintf("%f", body.Longitude),
		Latitude:    fmt.Sprintf("%f", body.Latitude),
	}

	// Save the record to the database
	if err := initializers.DB.Create(&securityRecord).Error; err != nil {
		c.JSON(500, gin.H{"message": "Failed to create security record"})
		return
	}

	// Prepare response
	response := struct {
		ID           uint   `json:"id"`
		SecurityId   uint   `json:"security_id"`
		SecurityName string `json:"security_name"`
		Block        string `json:"block"`
		PhoneNo      string `json:"phone_no"`
		Longitude    string `json:"longitude"`
		Latitude     string `json:"latitude"`
	}{
		ID:           securityRecord.ID,
		SecurityId:   securityID,
		SecurityName: securityUser.Name,
		Block:        securityRecord.Block,
		PhoneNo:      securityUser.Phone_No,
		Longitude:    securityRecord.Longitude,
		Latitude:     securityRecord.Latitude,
	}

	c.JSON(200, gin.H{
		"message": "Security record created successfully",
		"data":    response,
	})
}

func DeleteSecurityRecord(c *gin.Context) {
	if !isAdmin(c) {
		c.JSON(403, gin.H{"message": "Forbidden: Admins only"})
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

// isSecurity checks if the current user is a security
func isSecurity(c *gin.Context) bool {
	userID, exists := c.Get("user_id")
	if !exists {
		return false
	}
	uid, ok := userID.(uint)
	if !ok {
		return false
	}
	var user models.User
	if err := initializers.DB.First(&user, uid).Error; err != nil {
		return false
	}
	return user.Role_Id == 3 // 3 = Security
}

func GetSecurityRecordByDay(c *gin.Context) {
	if !userVerification(c) {
		c.JSON(403, gin.H{"message": "Forbidden: Users only"})
		return
	}

	// Get user id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"message": "User not authenticated"})
		return
	}
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(400, gin.H{"message": "Invalid user ID"})
		return
	}

	// Get date from query parameter
	date := c.Query("date")
	if date == "" {
		c.JSON(400, gin.H{"message": "Date query parameter is required (YYYY-MM-DD)"})
		return
	}

	// Parse the date and get start/end of day
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid date format. Use YYYY-MM-DD."})
		return
	}
	startOfDay := parsedDate
	endOfDay := parsedDate.Add(24 * time.Hour)

	// Get user's phone number
	var user models.User
	if err := initializers.DB.Select("phone_no").Where("id = ?", uid).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"message": "User not found"})
		return
	}

	type SecurityRecordResponse struct {
		ID           uint   `json:"id"`
		SecurityId   uint   `json:"security_id"`
		SecurityName string `json:"security_name"`
		Block        string `json:"block"`
		PhoneNo      string `json:"phone_no"`
		Longitude    string `json:"longitude"`
		Latitude     string `json:"latitude"`
	}

	var records []SecurityRecordResponse
	result := initializers.DB.Model(&models.SecurityRecord{}).
		Select("security_records.id, security_records.security_id, users.name as security_name, security_records.block, security_records.phone_no, security_records.longitude, security_records.latitude").
		Joins("left join users on users.id = security_records.security_id").
		Where("security_records.created_at >= ? AND security_records.created_at < ?", startOfDay, endOfDay).
		Scan(&records)
	if result.Error != nil {
		c.JSON(400, gin.H{"message": "Failed to get security records for the specified day"})
		return
	}

	c.JSON(200, gin.H{"data": records})
}

// GetSecurityRecordByUser returns all security records created by the authenticated security user
func GetSecurityRecordByUser(c *gin.Context) {
	if !isSecurity(c) {
		c.JSON(403, gin.H{"message": "Forbidden: Security only"})
		return
	}

	// Get security user id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"message": "User not authenticated"})
		return
	}
	securityID, ok := userID.(uint)
	if !ok {
		c.JSON(400, gin.H{"message": "Invalid user ID"})
		return
	}

	type SecurityRecordResponse struct {
		ID           uint   `json:"id"`
		SecurityId   uint   `json:"security_id"`
		SecurityName string `json:"security_name"`
		Block        string `json:"block"`
		PhoneNo      string `json:"phone_no"`
		Longitude    string `json:"longitude"`
		Latitude     string `json:"latitude"`
	}

	var records []SecurityRecordResponse
	result := initializers.DB.Model(&models.SecurityRecord{}).
		Select("security_records.id, security_records.security_id, users.name as security_name, security_records.block, security_records.phone_no, security_records.longitude, security_records.latitude").
		Joins("left join users on users.id = security_records.security_id").
		Where("security_records.security_id = ?", securityID).
		Scan(&records)
	if result.Error != nil {
		c.JSON(400, gin.H{"message": "Failed to get security records for this user"})
		return
	}

	c.JSON(200, gin.H{"data": records})
}

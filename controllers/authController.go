package controllers

import (
	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("86a044962f2416da2c15ebc88f2c9f828dc64f897c86720615a66a48cea37de5")

type Credentials struct {
	NIK      string `json:"nik"`
	Password string `json:"password"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {
	var creds Credentials

	if err := c.BindJSON(&creds); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	var user models.User

	if err := initializers.DB.Where("nik = ?", creds.NIK).First(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "User not found",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid password",
		})
		return
	}

	claims := &Claims{
		UserID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login success",
		"token":   tokenString,
	})
}

func Logout(c *gin.Context) {
	// Invalidate token

	c.JSON(200, gin.H{
		"message": "Logout success",
	})
}

func Authenticate(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{
			"message": "Invalid token",
		})
		return
	}

	c.Set("user_id", claims.UserID)
	c.Next()
}

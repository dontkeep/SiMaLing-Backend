package controllers

import (
	"time"

	"github.com/dontkeep/simaling-backend/initializers"
	"github.com/dontkeep/simaling-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("86a044962f2416da2c15ebc88f2c9f828dc64f897c86720615a66a48cea37de5")

type Credentials struct {
	Phone_No string `json:"phone_no"`
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

	if err := initializers.DB.Where("phone_no = ?", creds.Phone_No).First(&user).Error; err != nil {
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

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
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
		"role":    user.Role_Id,
	})
}

func Logout(c *gin.Context) {
	// Invalidate token
	tokenString, exists := c.Get("token")
	if !exists || tokenString == "" {
		c.JSON(400, gin.H{
			"message": "Token not found",
		})
		return
	}

	// Remove Bearer prefix before blacklisting
	rawToken := tokenString.(string)
	if len(rawToken) > 7 && rawToken[:7] == "Bearer " {
		rawToken = rawToken[7:]
	}

	if err := blacklistToken(rawToken); err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to blacklist token",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Logout success",
	})
}

func ExtractTokenMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	c.Set("token", tokenString)
	c.Next()
}

func Authenticate(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(400, gin.H{
			"message": "Token not found",
		})
		c.Abort()
		return
	}
	print("Token String: ", tokenString)
	rawToken := tokenString
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		rawToken = tokenString[7:]
	}
	// Check blacklist with full header (with Bearer)
	var blacklistToken models.BlacklistToken
	if err := initializers.DB.Where("token = ?", rawToken).First(&blacklistToken).Error; err == nil {
		c.JSON(401, gin.H{"message": "Token is blacklisted"})
		c.Abort()
		return
	}

	// Remove Bearer prefix for JWT parsing

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"message": "Invalid or expired token"})
		c.Abort()
		return
	}

	c.Set("user_id", claims.UserID)
	c.Set("token_string", tokenString)
	c.Next()
}

func blacklistToken(tokenString string) error {
	blacklistToken := models.BlacklistToken{Token: tokenString}
	result := initializers.DB.Create(&blacklistToken)
	return result.Error
}

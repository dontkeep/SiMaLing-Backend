package initializers

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseConnection() {
	var err error
	dsn := GetEnv("DB_USER", "root") + ":" + GetEnv("DB_PASSWORD", "1234") + "@tcp(" + GetEnv("DB_HOST", "mysql") + ":" + GetEnv("DB_PORT", "3306") + ")/" + GetEnv("DB_NAME", "simaling") + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database")
	}
}

func GetEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

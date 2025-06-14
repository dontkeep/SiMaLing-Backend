package initializers

import "os"

var AppConfig struct {
    BaseURL string
}

func LoadConfig() {
    AppConfig.BaseURL = os.Getenv("BASE_URL")
    // Default value if BASE_URL is not set

	// Set a value with real domain or IP address for production
	// or use localhost/local IP for development
    if AppConfig.BaseURL == "" {
        AppConfig.BaseURL = "http://localhost:3004"
    }
}
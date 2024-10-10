package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	AppID             string
	AppSecret         string
	APIURL            string
	AuthURL           string
	Port              string
	SingleChatUrl     string
	GroupChatUrl      string
	RegressionGroupID string
}

// LoadConfig loads the configuration from the .env file
func LoadConfig() *Config {
	err := godotenv.Load("/home/arvalinno/seatalk_bot/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		AppID:             os.Getenv("SEATALK_APP_ID"),
		AppSecret:         os.Getenv("SEATALK_APP_SECRET"),
		APIURL:            os.Getenv("SEATALK_API_URL"),
		AuthURL:           os.Getenv("SEATALK_AUTH_URL"),
		Port:              os.Getenv("PORT"),
		SingleChatUrl:     os.Getenv("SINGLE_CHAT_URL"),
		GroupChatUrl:      os.Getenv("SEATALK_SEND_GROUP_CHAT_URL"),
		RegressionGroupID: os.Getenv("REGRESSION_GROUP_ID"),
	}
}

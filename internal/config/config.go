package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
	"strings"
)

const Development = "development"
const Local = "local"
const DefaultPort = "8082"
const AppName = "ecom-api"

type Config struct {
	Env       string
	AppPort   string
	AppName   string
	JwtSecret string
	Database  struct {
		Host         string
		Port         int
		User         string
		Password     string
		DatabaseName string
		Schema       string
	}
}

var AppConfig = Config{
	AppPort: DefaultPort,
	Env:     Development,
	AppName: AppName,
}

func LoadConfig() {
	if envPort := os.Getenv("PORT"); envPort != "" {
		AppConfig.AppPort = envPort
	}

	if os.Getenv("APP_ENV") != "" {
		AppConfig.Env = strings.ToLower(os.Getenv("APP_ENV"))
	}

	if os.Getenv("APP_NAME") != "" {
		AppConfig.AppName = strings.ToLower(os.Getenv("APP_NAME"))
	}

	if os.Getenv("JWT_SECRET") != "" {
		AppConfig.JwtSecret = os.Getenv("JWT_SECRET")
	}

	setDatabaseConfig()
}

func setDatabaseConfig() {
	AppConfig.Database.User = os.Getenv("DB_USER")
	AppConfig.Database.Password = os.Getenv("DB_PASSWORD")
	AppConfig.Database.DatabaseName = os.Getenv("DB_DATABASE")
	AppConfig.Database.Schema = os.Getenv("DB_DEFAULT_SCHEMA")
	AppConfig.Database.Host = os.Getenv("DB_HOST")

	if os.Getenv("DB_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("DB_PORT"))

		if err != nil {
			panic("error parsing database port " + err.Error())
		}
		AppConfig.Database.Port = port
	}

}

func IsDev() bool {
	if AppConfig.Env == Development || AppConfig.Env == Local {
		return true
	}
	return false
}

func init() {
	LoadConfig()
}

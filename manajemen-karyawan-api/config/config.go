package config

import (
	"os"
)

var (
	DBUser       string
	DBPass       string
	DBHost       string
	DBPort       string
	DBName       string
	JWTSecret    string
	AppPort      string
	CookieDomain string
)

func InitConfig() {
	DBUser = os.Getenv("DB_USER")
	if DBUser == "" {
		DBUser = "root"
	}

	DBPass = os.Getenv("DB_PASS")

	DBHost = os.Getenv("DB_HOST")
	if DBHost == "" {
		DBHost = "127.0.0.1"
	}

	DBPort = os.Getenv("DB_PORT")
	if DBPort == "" {
		DBPort = "3306"
	}

	DBName = os.Getenv("DB_NAME")
	if DBName == "" {
		DBName = "manajemen_karyawan"
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		JWTSecret = "default_jwt_secret"
	}

	AppPort = os.Getenv("APP_PORT")
	if AppPort == "" {
		AppPort = "8080"
	}

	CookieDomain = os.Getenv("COOKIE_DOMAIN")
	if CookieDomain == "" {
		CookieDomain = "localhost"
	}
}

package config

import (
	"github.com/wiormiw/GrowBaks/helper"
)

type (
	// Configuration is application configuration
	Configuration struct {
		Const    *Constants
		Database *Database
		Secret   *Secret
	}

	// Constants is used To store configurable value, not the constant-constant value
	Constants struct {
		HTTPPort      int
		ENV           string
		CloudinaryUrl string
	}

	// Database configuration
	Database struct {
		DBUser     string
		DBPassword string
		DBHost     string
		DBPort     int
		DBName     string
	}

	// Secret configuration
	Secret struct {
		JWTSecret string
	}
)

// LoadConfiguration...
func LoadConfiguration() *Configuration {
	return &Configuration{
		Const:    loadConstants(),
		Database: loadDatabase(),
		Secret:   loadSecret(),
	}
}

// loadConstants...
func loadConstants() *Constants {
	return &Constants{
		HTTPPort:      helper.GetEnvInt("PORT"),
		ENV:           helper.GetEnvString("ENV"),
		CloudinaryUrl: helper.GetEnvString("CLOUDINARY_URL"),
	}
}

// loadDatabase...
func loadDatabase() *Database {
	return &Database{
		DBHost:     helper.GetEnvString("DB_HOST"),
		DBUser:     helper.GetEnvString("DB_USER"),
		DBPassword: helper.GetEnvString("DB_PASSWORD"),
		DBPort:     helper.GetEnvInt("DB_PORT"),
		DBName:     helper.GetEnvString("DB_NAME"),
	}
}

// loadSecret...
func loadSecret() *Secret {
	return &Secret{
		JWTSecret: helper.GetEnvString("JWT_SECRET"),
	}
}

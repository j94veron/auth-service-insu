package config

import (
	"encoding/base64"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

// ConnectDB declares and returns the database connection based on the environment
func ConnectDB() (*gorm.DB, error) {
	env := os.Getenv("CONFIG_ENV")
	var dsnBase64 string

	switch env {
	case "PROD":
		dsnBase64 = os.Getenv("DB_DSN_PROD")
	case "TEST":
		dsnBase64 = os.Getenv("DB_DSN_TEST")
	default:
		return nil, fmt.Errorf("unknown environment inCONFIG_ENV: %s", env)
	}

	if dsnBase64 == "" {
		return nil, fmt.Errorf("environment variable for base64 DSN not defined for environment %s", env)
	}

	// Decode base64
	dsnBytes, err := base64.StdEncoding.DecodeString(dsnBase64)
	if err != nil {
		return nil, fmt.Errorf("error decoding base6: %v", err)
	}
	dsn := string(dsnBytes)

	//Add the connection to the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

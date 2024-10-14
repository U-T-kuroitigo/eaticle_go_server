package configuration

import (
	"log"
	"os"

	"fmt"

	"github.com/U-T-kuroitigo/eaticle_go_server/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Configuration creates a struct for the environment variables
type Configuration struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
}

// GetConfiguration gets the configuration from the environment variables
func GetConfiguration() Configuration {
	var c Configuration
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 環境変数を取得
	c.Server = os.Getenv("Server")
	c.Port = os.Getenv("Port")
	c.User = os.Getenv("User")
	c.Password = os.Getenv("Password")
	c.Database = os.Getenv("Database")

	return c
}

// InitDB initializes the database connection and performs auto-migration
func InitDB() *gorm.DB {
	c := GetConfiguration()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Server, c.Port, c.User, c.Password, c.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// AutoMigrate the models
	db.AutoMigrate(&models.User{})

	return db
}
